package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/config"
	"github.com/yuwen002/go-meteor-cms/ent"
	"github.com/yuwen002/go-meteor-cms/ent/adminpermission"
	"github.com/yuwen002/go-meteor-cms/ent/adminrole"
	"github.com/yuwen002/go-meteor-cms/ent/adminuser"
	"github.com/yuwen002/go-meteor-cms/internal/common"
	"github.com/yuwen002/go-meteor-cms/internal/utils"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var permissionBypassRoutes = map[string]struct{}{
	http.MethodPost + " /admin/logout":                        {},
	http.MethodGet + " /admin/test-token":                     {},
	http.MethodPut + " /admin/admin-users/me/change-password": {},
	// 个人中心相关接口必须放行：普通管理员在拿到任何权限之前需要取回菜单，
	// 并需要维护自己的资料（不涉及角色与启用状态）
	http.MethodGet + " /admin/profile/permissions": {},
	http.MethodGet + " /admin/profile":             {},
	http.MethodPut + " /admin/profile":             {},
}

type PermissionMiddleware struct {
	config    *config.Config
	db        *ent.Client
	permCache *collection.Cache
}

type matchedPermission struct {
	permission *ent.AdminPermission
	staticHits int
}

func NewPermissionMiddleware(c *config.Config, db *ent.Client) rest.Middleware {
	cache, err := collection.NewCache(60 * time.Second)
	if err != nil {
		panic(err)
	}
	m := &PermissionMiddleware{
		config:    c,
		db:        db,
		permCache: cache,
	}
	return m.Handle
}

func (m *PermissionMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := utils.CurrentUserID(r.Context())
		if !ok {
			common.Fail(w, common.ErrMissingUserContext, common.GetErrorMessage(common.ErrMissingUserContext))
			return
		}

		user, err := m.db.AdminUser.
			Query().
			Where(adminuser.IDEQ(userID)).
			Only(r.Context())

		if err != nil {
			if ent.IsNotFound(err) {
				common.Fail(w, common.ErrMissingUserContext, common.GetErrorMessage(common.ErrMissingUserContext))
				return
			}
			logx.Errorf("query admin user failed: %v", err)
			common.Fail(w, common.ErrInternalServer, common.GetErrorMessage(common.ErrInternalServer))
			return
		}

		if !user.IsActive {
			common.Fail(w, common.ErrAccountInactive, common.GetErrorMessage(common.ErrAccountInactive))
			return
		}

		if user.IsSuper || isPermissionBypassRoute(r.Method, r.URL.Path) {
			next(w, r)
			return
		}

		permission, err := m.findPermissionForRequest(r.Context(), r.Method, r.URL.Path)
		if err != nil {
			logx.Errorf("match api permission failed: method=%s path=%s err=%v", r.Method, r.URL.Path, err)
			common.Fail(w, common.ErrInternalServer, common.GetErrorMessage(common.ErrInternalServer))
			return
		}
		if permission == nil {
			logx.Errorf("api permission missing: user_id=%d method=%s path=%s", userID, r.Method, r.URL.Path)
			common.Fail(w, common.ErrPermissionDenied, common.GetErrorMessage(common.ErrPermissionDenied))
			return
		}

		hasPermission, err := m.checkUserPermission(r.Context(), userID, permission.Permission)
		if err != nil {
			logx.Errorf("check user permission failed: user_id=%d permission=%s err=%v", userID, permission.Permission, err)
			common.Fail(w, common.ErrInternalServer, common.GetErrorMessage(common.ErrInternalServer))
			return
		}
		if !hasPermission {
			logx.Errorf("permission denied: user_id=%d method=%s path=%s permission=%s", userID, r.Method, r.URL.Path, permission.Permission)
			common.Fail(w, common.ErrPermissionDenied, common.GetErrorMessage(common.ErrPermissionDenied))
			return
		}

		next(w, r)
	}
}

func (m *PermissionMiddleware) checkUserPermission(ctx context.Context, userID int64, permission string) (bool, error) {
	perms, err := m.loadUserPermissions(ctx, userID)
	if err != nil {
		return false, err
	}
	for _, p := range perms {
		if p == permission {
			return true, nil
		}
	}
	return false, nil
}

// loadUserPermissions 加载用户拥有的全部权限标识，结果按 userID 缓存（短 TTL），
// 避免每次请求都走「用户→角色→权限」多表查询。注意：用户记录（is_super / is_active）
// 仍每请求查库，以保证权限变更实时生效（见任务1），此处只缓存「角色→权限」映射。
func (m *PermissionMiddleware) loadUserPermissions(ctx context.Context, userID int64) ([]string, error) {
	key := "perm:" + strconv.FormatInt(userID, 10)
	v, err := m.permCache.Take(key, func() (interface{}, error) {
		var perms []string
		err := m.db.AdminUser.
			Query().
			Where(adminuser.IDEQ(userID)).
			QueryRoles().
			Where(adminrole.IsActive(true)).
			QueryPermissions().
			Where(adminpermission.IsActive(true)).
			Select(adminpermission.FieldPermission).
			Scan(ctx, &perms)
		if err != nil {
			return nil, err
		}
		return perms, nil
	})
	if err != nil {
		return nil, err
	}
	return v.([]string), nil
}

func (m *PermissionMiddleware) findPermissionForRequest(ctx context.Context, method, requestPath string) (*ent.AdminPermission, error) {
	key := "apiperm:" + strings.ToUpper(method) + ":" + normalizePath(requestPath)
	v, err := m.permCache.Take(key, func() (interface{}, error) {
		return m.findPermissionForRequestUncached(ctx, method, requestPath)
	})
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	p, ok := v.(*ent.AdminPermission)
	if !ok || p == nil {
		return nil, nil
	}
	return p, nil
}

func (m *PermissionMiddleware) findPermissionForRequestUncached(ctx context.Context, method, requestPath string) (*ent.AdminPermission, error) {
	permissions, err := m.db.AdminPermission.
		Query().
		Where(
			adminpermission.TypeEQ(3),
			adminpermission.IsActive(true),
			adminpermission.MethodEqualFold(strings.ToUpper(method)),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}

	var best *matchedPermission
	for _, permission := range permissions {
		staticHits, matched := matchPermissionPath(permission.APIPath, requestPath)
		if !matched {
			continue
		}

		candidate := &matchedPermission{
			permission: permission,
			staticHits: staticHits,
		}

		if best == nil || candidate.staticHits > best.staticHits {
			best = candidate
			continue
		}

		if candidate.staticHits == best.staticHits {
			logx.Errorf("ambiguous api permission config: method=%s path=%s matched %q and %q", method, requestPath, best.permission.APIPath, candidate.permission.APIPath)
			return nil, nil
		}
	}

	if best == nil {
		return nil, nil
	}

	return best.permission, nil
}

func isPermissionBypassRoute(method, path string) bool {
	_, ok := permissionBypassRoutes[strings.ToUpper(method)+" "+normalizePath(path)]
	return ok
}

func matchPermissionPath(templatePath, requestPath string) (int, bool) {
	templateSegments := splitPathSegments(templatePath)
	requestSegments := splitPathSegments(requestPath)
	if len(templateSegments) != len(requestSegments) {
		return 0, false
	}

	staticHits := 0
	for i := range templateSegments {
		templateSegment := templateSegments[i]
		requestSegment := requestSegments[i]

		if isPathParamSegment(templateSegment) {
			if requestSegment == "" {
				return 0, false
			}
			continue
		}

		if templateSegment != requestSegment {
			return 0, false
		}
		staticHits++
	}

	return staticHits, true
}

func splitPathSegments(path string) []string {
	normalized := normalizePath(path)
	if normalized == "/" {
		return []string{}
	}

	return strings.Split(strings.TrimPrefix(normalized, "/"), "/")
}

func normalizePath(path string) string {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return "/"
	}

	trimmed = strings.TrimSuffix(trimmed, "/")
	if trimmed == "" {
		return "/"
	}
	if !strings.HasPrefix(trimmed, "/") {
		return "/" + trimmed
	}

	return trimmed
}

func isPathParamSegment(segment string) bool {
	return strings.HasPrefix(segment, ":") && len(segment) > 1
}
