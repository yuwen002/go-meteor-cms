package middleware

import (
	"context"
	"math"
	"net/http"
	"strings"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/config"
	"github.com/yuwen002/go-meteor-cms/ent"
	"github.com/yuwen002/go-meteor-cms/ent/adminpermission"
	"github.com/yuwen002/go-meteor-cms/ent/adminrole"
	"github.com/yuwen002/go-meteor-cms/ent/adminuser"
	"github.com/yuwen002/go-meteor-cms/internal/common"
	"github.com/yuwen002/go-meteor-cms/internal/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var permissionBypassRoutes = map[string]struct{}{
	http.MethodPost + " /admin/logout":                        {},
	http.MethodGet + " /admin/test-token":                     {},
	http.MethodPut + " /admin/admin-users/me/change-password": {},
}

type PermissionMiddleware struct {
	config *config.Config
	db     *ent.Client
}

type matchedPermission struct {
	permission *ent.AdminPermission
	staticHits int
}

func NewPermissionMiddleware(c *config.Config, db *ent.Client) rest.Middleware {
	m := &PermissionMiddleware{
		config: c,
		db:     db,
	}
	return m.Handle
}

func (m *PermissionMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userClaims := utils.GetUserFromCtx(r.Context())
		if userClaims == nil {
			common.Fail(w, common.ErrUnauthorized, common.GetErrorMessage(common.ErrUnauthorized))
			return
		}

		userID, ok := extractUserID(userClaims)
		if !ok {
			common.Fail(w, common.ErrUnauthorized, common.GetErrorMessage(common.ErrUnauthorized))
			return
		}

		user, err := m.db.AdminUser.
			Query().
			Where(adminuser.IDEQ(userID)).
			Only(r.Context())

		if err != nil {
			if ent.IsNotFound(err) {
				common.Fail(w, common.ErrUnauthorized, common.GetErrorMessage(common.ErrUnauthorized))
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
	count, err := m.db.AdminUser.
		Query().
		Where(adminuser.IDEQ(userID)).
		QueryRoles().
		Where(adminrole.IsActive(true)).
		QueryPermissions().
		Where(
			adminpermission.IsActive(true),
			adminpermission.PermissionEQ(permission),
		).
		Count(ctx)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (m *PermissionMiddleware) findPermissionForRequest(ctx context.Context, method, requestPath string) (*ent.AdminPermission, error) {
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

func extractUserID(claims map[string]interface{}) (int64, bool) {
	value, ok := claims["user_id"]
	if !ok {
		return 0, false
	}

	switch v := value.(type) {
	case int64:
		return v, true
	case int:
		return int64(v), true
	case float64:
		if math.Trunc(v) != v {
			return 0, false
		}
		return int64(v), true
	default:
		return 0, false
	}
}
