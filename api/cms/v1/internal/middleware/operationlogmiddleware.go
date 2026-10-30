package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/config"
	"github.com/yuwen002/go-meteor-cms/ent"
	"github.com/yuwen002/go-meteor-cms/internal/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

const (
	// maxLogBodyBytes 请求参数快照的最大长度，超出部分截断
	maxLogBodyBytes = 4096
	// writeLogTimeout 写日志的超时时间，避免拖慢请求
	writeLogTimeout = 3 * time.Second
)

// sensitiveKeys 需要脱敏的参数名，避免明文密码进入日志
var sensitiveKeys = []string{
	"password", "old_password", "new_password", "confirm_password", "captcha", "token",
}

// moduleByPrefix 请求路径前缀到业务模块的映射，用于日志归类
var moduleByPrefix = map[string]string{
	"/admin/admin-users": "管理员管理",
	"/admin/departments": "部门管理",
	"/admin/roles":       "角色管理",
	"/admin/permissions": "权限管理",
	"/admin/profile":     "个人中心",
	"/admin/logs":        "日志管理",
}

// OperationLogMiddleware 操作日志中间件
//
// 记录受保护接口的调用情况：操作人、方法、路径、来源 IP、结果、耗时与脱敏后的参数快照。
// 日志写入放在独立 goroutine 中，不阻塞请求；file upload 等二进制请求不记录参数。
type OperationLogMiddleware struct {
	config config.Config
	db     *ent.Client
}

// NewOperationLogMiddleware 创建操作日志中间件
func NewOperationLogMiddleware(c config.Config, db *ent.Client) rest.Middleware {
	m := &OperationLogMiddleware{config: c, db: db}
	return m.Handle
}

// Handle 返回中间件处理函数
func (m *OperationLogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 读取请求体快照并回填，保证后续 handler 仍能正常解析
		body := readBodySnapshot(r)

		rec := &responseRecorder{ResponseWriter: w, status: http.StatusOK, body: bytes.NewBuffer(nil)}
		next(rec, r)

		duration := time.Since(start)
		go m.save(r, rec, body, duration)
	}
}

// save 异步写入操作日志
func (m *OperationLogMiddleware) save(r *http.Request, rec *responseRecorder, body string, duration time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), writeLogTimeout)
	defer cancel()

	var userID int64
	var username string
	if id, ok := utils.CurrentUserID(r.Context()); ok {
		userID = id
	}
	if claims := utils.GetUserFromCtx(r.Context()); claims != nil {
		if name, ok := claims["username"].(string); ok {
			username = name
		}
	}

	code, msg := rec.businessResult()

	create := m.db.OperationLog.Create().
		SetUserID(userID).
		SetUsername(username).
		SetModule(resolveModule(r.URL.Path)).
		SetMethod(r.Method).
		SetPath(r.URL.Path).
		SetIP(utils.ClientIPFromRequest(r)).
		SetUserAgent(utils.Truncate(r.UserAgent(), 512)).
		SetDurationMs(duration.Milliseconds())

	if body != "" {
		create = create.SetParams(body)
	}

	if code == 0 {
		create = create.SetStatus(0)
	} else {
		create = create.SetStatus(1).SetErrorMsg(utils.Truncate(msg, 512))
	}

	if _, err := create.Save(ctx); err != nil {
		logx.Errorf("写入操作日志失败: %v", err)
	}
}

// readBodySnapshot 读取请求体快照。
//
// 只读取有 body 且非二进制的请求，读取后原样回填 r.Body；
// 超长内容会被截断，敏感字段统一替换为 ******。
func readBodySnapshot(r *http.Request) string {
	if r.Body == nil || r.Method == http.MethodGet || r.Method == http.MethodHead {
		return ""
	}

	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/") {
		// 文件上传等二进制请求不记录内容，避免日志膨胀
		return ""
	}

	raw, err := io.ReadAll(io.LimitReader(r.Body, maxLogBodyBytes))
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return ""
	}
	// 回填，保证下游 handler 能继续读取
	r.Body = io.NopCloser(bytes.NewReader(raw))

	snapshot := strings.TrimSpace(string(raw))
	if snapshot == "" {
		return ""
	}

	return sanitizeParams(contentType, snapshot)
}

// sanitizeParams 脱敏参数中的敏感字段
func sanitizeParams(contentType, snapshot string) string {
	if !strings.Contains(contentType, "application/json") {
		// 非 JSON 表单无法可靠逐字段脱敏，直接丢弃以免泄露敏感值
		return ""
	}

	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(snapshot), &payload); err != nil {
		return ""
	}

	for key := range payload {
		for _, sensitive := range sensitiveKeys {
			if strings.EqualFold(key, sensitive) {
				payload[key] = "******"
				break
			}
		}
	}

	masked, err := json.Marshal(payload)
	if err != nil {
		return ""
	}

	return utils.Truncate(string(masked), maxLogBodyBytes)
}

// resolveModule 按路径前缀归类业务模块
func resolveModule(path string) string {
	for prefix, module := range moduleByPrefix {
		if strings.HasPrefix(path, prefix) {
			return module
		}
	}
	return "其他"
}

// responseRecorder 捕获响应状态码与响应体，用于判定操作结果
type responseRecorder struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

// WriteHeader 记录状态码
func (r *responseRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// Write 同时写入底层响应与缓冲
func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// businessResult 解析响应体中的业务状态码与提示信息
func (r *responseRecorder) businessResult() (int, string) {
	var out struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.Unmarshal(r.body.Bytes(), &out); err != nil {
		// 非 JSON 响应（如静态资源）按 HTTP 状态码判定
		if r.status >= http.StatusBadRequest {
			return r.status, http.StatusText(r.status)
		}
		return 0, ""
	}
	return out.Code, out.Msg
}
