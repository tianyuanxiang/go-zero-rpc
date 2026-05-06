// Package response 提供统一的HTTP响应结构和辅助函数。
// 所有接口返回数据均通过本包封装，保证响应格式一致性。
package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 业务响应码常量定义
const (
	// CodeSuccess 请求成功
	CodeSuccess = 0
	// CodeParamError 参数错误
	CodeParamError = 10001
	// CodeUnauthorized 未授权（未登录或token无效）
	CodeUnauthorized = 10002
	// CodeForbidden 权限不足
	CodeForbidden = 10003
	// CodeNotFound 数据不存在
	CodeNotFound = 10004
	// CodeDuplicate 数据重复/已存在
	CodeDuplicate = 10005
	// CodeInternalError 服务内部错误
	CodeInternalError = 10006
	// CodeTokenExpired token已过期
	CodeTokenExpired = 10007
	// CodeAccountDisabled 账号已禁用
	CodeAccountDisabled = 10008
	// CodePasswordError 密码错误
	CodePasswordError = 10009
)

// 业务响应消息常量定义
const (
	MsgSuccess         = "success"
	MsgParamError      = "参数错误"
	MsgUnauthorized    = "未授权，请先登录"
	MsgForbidden       = "权限不足，无法访问"
	MsgNotFound        = "数据不存在"
	MsgDuplicate       = "数据已存在"
	MsgInternalError   = "服务器开小差了，请联系系统运维人员"
	MsgTokenExpired    = "登录已过期，请重新登录"
	MsgAccountDisabled = "账号已禁用"
	MsgPasswordError   = "用户名或密码错误"
)

// Response 统一HTTP响应结构体。
// 所有接口返回的JSON数据均符合该结构。
type Response struct {
	// Code 业务响应码，0表示成功，非0表示失败
	Code int `json:"code"`
	// Msg 响应消息，成功时为"success"，失败时为错误描述
	Msg string `json:"msg"`
	// Data 响应数据，成功时携带业务数据，失败时可为null
	Data interface{} `json:"data,omitempty"`
}

// OK 返回成功响应（无数据）。
// 参数 w 为 http.ResponseWriter，用于写入HTTP响应。
// 参数 r 为 *http.Request，用于读取请求上下文。
func OK(w http.ResponseWriter, r *http.Request) {
	httpx.OkJsonCtx(r.Context(), w, Response{
		Code: CodeSuccess,
		Msg:  MsgSuccess,
		Data: nil,
	})
}

// OkWithData 返回成功响应（携带业务数据）。
// 参数 data 为需要返回给前端的业务数据。
func OkWithData(w http.ResponseWriter, r *http.Request, data interface{}) {
	httpx.OkJsonCtx(r.Context(), w, Response{
		Code: CodeSuccess,
		Msg:  MsgSuccess,
		Data: data,
	})
}

// Fail 返回失败响应（使用预定义错误码和消息）。
// 参数 code 为业务错误码，参数 msg 为错误描述信息。
func Fail(w http.ResponseWriter, r *http.Request, code int, msg string) {
	httpx.OkJsonCtx(r.Context(), w, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

// FailWithMsg 返回失败响应（使用参数错误码，自定义消息）。
// 常用于参数校验失败时快速返回。
func FailWithMsg(w http.ResponseWriter, r *http.Request, msg string) {
	httpx.OkJsonCtx(r.Context(), w, Response{
		Code: CodeParamError,
		Msg:  msg,
		Data: nil,
	})
}

// FailParam 返回参数错误响应。
func FailParam(w http.ResponseWriter, r *http.Request) {
	Fail(w, r, CodeParamError, MsgParamError)
}

// FailUnauthorized 返回未授权响应。
func FailUnauthorized(w http.ResponseWriter, r *http.Request) {
	Fail(w, r, CodeUnauthorized, MsgUnauthorized)
}

// FailForbidden 返回权限不足响应。
func FailForbidden(w http.ResponseWriter, r *http.Request) {
	Fail(w, r, CodeForbidden, MsgForbidden)
}

// FailNotFound 返回数据不存在响应。
func FailNotFound(w http.ResponseWriter, r *http.Request) {
	Fail(w, r, CodeNotFound, MsgNotFound)
}

// FailInternal 返回服务内部错误响应。
func FailInternal(w http.ResponseWriter, r *http.Request) {
	Fail(w, r, CodeInternalError, MsgInternalError)
}
