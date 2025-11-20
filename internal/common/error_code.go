package common

const (
	// Success 成功
	Success = 0

	// Common errors (1000-1999)
	ErrInternalServer = 1000 // 服务器内部错误

	// Auth errors (2000-2999)
	ErrUnauthorized = 2001 // 未授权
	ErrForbidden    = 2003 // 禁止访问

	// Admin user errors (3000-3999)
	ErrAdminUserAlreadyExists = 3001 // 管理员已存在
	ErrAdminUserNotFound      = 3002 // 管理员不存在
	ErrAdminPasswordIncorrect = 3003 // 密码错误

	// Request validation errors (4000-4999)
	ErrInvalidParams = 4001 // 参数错误
)

var errorMessages = map[int]string{
	Success:                   "success",
	ErrInternalServer:         "服务器内部错误",
	ErrUnauthorized:           "未授权",
	ErrForbidden:              "禁止访问",
	ErrAdminUserAlreadyExists: "管理员已存在",
	ErrAdminUserNotFound:      "管理员不存在",
	ErrAdminPasswordIncorrect: "密码错误",
	ErrInvalidParams:          "参数错误",
}
