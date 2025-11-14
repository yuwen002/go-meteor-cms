package common

type BizError struct {
	Code int
	Msg  string
}

func NewBizError(code int, msg string) *BizError {
	return &BizError{Code: code, Msg: msg}
}

func (e *BizError) Error() string {
	return e.Msg
}
