package common

// BizError 业务错误结构体
type BizError struct {
	Code int    // 错误码
	Msg  string // 错误信息
}

// Error 实现 error 接口，用于将 BizError 作为标准 error 使用
// 这是 Go 标准库中 error 接口要求的实现
func (e *BizError) Error() string {
	return e.Msg
}

// NewBizError 创建一个新的业务错误
// 使用 errorMessages 中定义的错误信息
//
// 参数:
//   - code: 错误码，对应 errorMessages 中的键
//
// 返回值:
//   - *BizError: 返回业务错误实例
func NewBizError(code int) *BizError {
	msg, exists := errorMessages[code]
	if !exists {
		msg = "未知错误"
	}
	return &BizError{Code: code, Msg: msg}
}

// NewBizErrorWithMsg 创建一个带自定义错误信息的业务错误
//
// 参数:
//   - code: 错误码
//   - msg: 自定义错误信息
//
// 返回值:
//   - *BizError: 返回业务错误实例
func NewBizErrorWithMsg(code int, msg string) *BizError {
	return &BizError{
		Code: code,
		Msg:  msg,
	}
}
