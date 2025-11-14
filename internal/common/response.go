package common

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// Response 统一响应结构体
// 用于规范API的返回格式，包含状态码、消息和数据
// 状态码说明：
// 0   - 成功
// >0  - 已知业务错误（如：参数错误、权限不足等）
// <0  - 系统错误（如：数据库错误、第三方服务错误等）
type Response struct {
	Code int         `json:"code"` // 状态码
	Msg  string      `json:"msg"`  // 提示信息
	Data interface{} `json:"data,omitempty"` // 返回数据
}

// Ok 返回成功响应
// 参数：
//   - w: HTTP响应写入器
//   - data: 需要返回的数据，可以是任意类型
//
// 示例：
//   common.Ok(w, userInfo)
func Ok(w http.ResponseWriter, data interface{}) {
	httpx.OkJson(w, Response{
		Code: 0,           // 0 表示成功
		Msg:  "success",   // 成功提示信息
		Data: data,        // 返回的数据
	})
}

// Fail 返回失败响应
// 参数：
//   - w: HTTP响应写入器
//   - code: 错误码，非0表示错误
//   - msg: 错误信息
//
// 示例：
//   common.Fail(w, 1001, "用户不存在")
func Fail(w http.ResponseWriter, code int, msg string) {
	httpx.OkJson(w, Response{
		Code: code,    // 错误码
		Msg:  msg,     // 错误信息
		Data: nil,     // 错误时一般不返回数据
	})
}
