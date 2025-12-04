package common

const (
	// Success 成功
	Success = 0

	// Common errors (1000-1999)
	ErrInternalServer = 1000 // 服务器内部错误

	// Auth errors (2000-2999)
	ErrUnauthorized       = 2001 // 未授权
	ErrForbidden          = 2003 // 禁止访问
	ErrMissingUserContext = 2004 // 用户上下文信息缺失
	ErrMissingCaptcha     = 2005 // 验证码不能为空
	ErrInvalidCaptcha     = 2006 // 验证码错误
	ErrAccountInactive    = 2007 // 账号未激活
	ErrCaptchaGenerate    = 2008 // 验证码生成失败
	ErrInvalidToken       = 2009 // 无效的token

	// Admin user errors (3000-3999)
	ErrAdminUserAlreadyExists  = 3001 // 管理员已存在
	ErrAdminUserNotFound       = 3002 // 管理员不存在
	ErrAdminPasswordIncorrect  = 3003 // 密码错误
	ErrAdminCountFailed        = 3004 // 获取管理员总数失败
	ErrAdminListFailed         = 3005 // 获取管理员列表失败
	ErrAdminPasswordHashFail   = 3006 // 密码加密失败
	ErrAdminCreateFailed       = 3007 // 创建管理员失败
	ErrAdminEmailAlreadyExists = 3008 // 邮箱已被注册
	ErrUserIDFormat            = 3009 // 用户ID格式错误
	ErrPasswordUpdateFailed    = 3010 // 密码更新失败
	ErrCannotResetOwnPassword  = 3011 // 不能重置自己的密码

	// Request validation errors (4000-4999)
	ErrInvalidParams = 4001 // 参数错误

	// Department errors (5000-5999)
	ErrDepartmentNotFound      = 5001 // 部门不存在
	ErrDepartmentCreateFail    = 5002 // 创建部门失败
	ErrDepartmentUpdateFail    = 5003 // 更新部门失败
	ErrDepartmentDeleteFail    = 5004 // 删除部门失败
	ErrDepartmentListFail      = 5005 // 获取部门列表失败
	ErrDepartmentParentNotExist = 5006 // 父级部门不存在
)

var errorMessages = map[int]string{
	Success:                    "success",
	ErrInternalServer:          "服务器内部错误",
	ErrUnauthorized:            "用户名或密码错误",
	ErrForbidden:               "禁止访问",
	ErrMissingUserContext:      "用户会话已过期，请重新登录",
	ErrInvalidToken:            "无效的token",
	ErrMissingCaptcha:          "验证码不能为空",
	ErrInvalidCaptcha:          "验证码错误",
	ErrAccountInactive:         "账号未激活，请联系管理员",
	ErrCaptchaGenerate:         "验证码生成失败，请重试",
	ErrAdminUserAlreadyExists:  "管理员已存在",
	ErrAdminUserNotFound:       "管理员不存在",
	ErrAdminPasswordIncorrect:  "密码错误",
	ErrAdminCountFailed:        "获取管理员总数失败",
	ErrAdminListFailed:         "获取管理员列表失败",
	ErrAdminPasswordHashFail:   "密码加密失败",
	ErrAdminCreateFailed:       "创建管理员失败",
	ErrAdminEmailAlreadyExists: "邮箱已被注册",
	ErrUserIDFormat:            "用户ID格式错误",
	ErrPasswordUpdateFailed:    "密码更新失败",
	ErrCannotResetOwnPassword:  "不能重置自己的密码，请使用修改密码功能",
	ErrInvalidParams:           "参数错误",
	// Department error messages
	ErrDepartmentNotFound:      "部门不存在",
	ErrDepartmentCreateFail:    "创建部门失败",
	ErrDepartmentUpdateFail:    "更新部门失败",
	ErrDepartmentDeleteFail:    "删除部门失败",
	ErrDepartmentListFail:      "获取部门列表失败",
	ErrDepartmentParentNotExist: "父级部门不存在",
}

// GetErrorMessage 获取错误码对应的错误信息
// 参数：
//   - code: 错误码
//
// 返回值：
//   - string: 错误信息，如果错误码不存在则返回"未知错误"
func GetErrorMessage(code int) string {
	if msg, ok := errorMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
