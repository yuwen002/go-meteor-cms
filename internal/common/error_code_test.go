package common

import "testing"

// TestErrorMessageContracts 固化各鉴权场景的错误文案，避免语义串味。
//
// 历史上 ErrUnauthorized 的文案是「用户名或密码错误」，
// 导致「缺少凭证」「会话过期」等场景也提示密码错误，具有误导性。
func TestErrorMessageContracts(t *testing.T) {
	cases := map[int]string{
		Success:               "success",
		ErrUnauthorized:       "请先登录",
		ErrMissingUserContext: "用户会话已过期，请重新登录",
		ErrInvalidToken:       "无效的token",
		ErrLoginFailed:        "用户名或密码错误",
		ErrPermissionDenied:   "权限不足",
		ErrAccountInactive:    "账号未激活，请联系管理员",
	}

	for code, want := range cases {
		if got := GetErrorMessage(code); got != want {
			t.Errorf("GetErrorMessage(%d) = %q, want %q", code, got, want)
		}
	}
}

// TestLoginFailedIsDistinctFromUnauthorized 登录失败与未登录必须使用不同错误码
func TestLoginFailedIsDistinctFromUnauthorized(t *testing.T) {
	if ErrLoginFailed == ErrUnauthorized {
		t.Fatal("登录失败禁止复用未授权错误码")
	}
	if GetErrorMessage(ErrLoginFailed) == GetErrorMessage(ErrUnauthorized) {
		t.Fatal("登录失败与未授权的提示文案不应相同")
	}
}

func TestGetErrorMessageUnknownCode(t *testing.T) {
	if got := GetErrorMessage(999999); got != "未知错误" {
		t.Fatalf("未知错误码应返回兜底文案, got %q", got)
	}
}

func TestNewBizError(t *testing.T) {
	// 已注册的错误码使用库内文案
	err := NewBizError(ErrLoginFailed)
	if err.Code != ErrLoginFailed || err.Msg != "用户名或密码错误" {
		t.Fatalf("NewBizError = %+v", err)
	}
	if err.Error() != "用户名或密码错误" {
		t.Fatalf("Error() 应返回文案, got %q", err.Error())
	}

	// 未注册的错误码回落到兜底文案
	unknown := NewBizError(999999)
	if unknown.Msg != "未知错误" {
		t.Fatalf("未注册错误码应回落兜底文案, got %q", unknown.Msg)
	}

	custom := NewBizErrorWithMsg(ErrAdminUserNotFound, "目标管理员不存在")
	if custom.Code != ErrAdminUserNotFound || custom.Msg != "目标管理员不存在" {
		t.Fatalf("NewBizErrorWithMsg = %+v", custom)
	}
}
