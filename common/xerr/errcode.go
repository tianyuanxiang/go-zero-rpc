package xerr

import (
	"fmt"
	"strconv"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BusinessErrorReason ErrorInfo.Reason 字段固定值，
// 客户端据此识别一个 gRPC status 是否承载本项目的业务错误码。
const BusinessErrorReason = "BUSINESS_ERROR"

// BusinessErrorDomain ErrorInfo.Domain 字段固定值，标识错误所属域。
const BusinessErrorDomain = "go-zero-rpc"

// 业务错误码定义
// 编码规则：
//   - 0      : 成功
//   - 10001-10099 : 通用系统错误
//   - 10100-10199 : 认证授权相关错误
//   - 10200-10299 : 用户管理相关错误
//   - 10300-10399 : 角色菜单相关错误
//   - 10400-10499 : 业务数据相关错误（槽体、事件等）
const (
	// ErrSuccess 成功
	ErrSuccess = 0

	// --- 通用系统错误 ---

	// ErrParamInvalid 参数校验失败
	ErrParamInvalid = 10001
	// ErrUnauthorized 未登录或Token无效
	ErrUnauthorized = 10002
	// ErrForbidden 权限不足
	ErrForbidden = 10003
	// ErrNotFound 目标数据不存在
	ErrNotFound = 10004
	// ErrDuplicate 数据重复，违反唯一性约束
	ErrDuplicate = 10005
	// ErrInternal 服务内部错误
	ErrInternal = 10006
	// ErrTokenExpired Token已过期
	ErrTokenExpired = 10007
	// ErrTokenInvalid Token格式无效
	ErrTokenInvalid = 10008

	// --- 认证授权相关错误 ---

	// ErrAccountDisabled 账号已被禁用
	ErrAccountDisabled = 10100
	// ErrPasswordWrong 用户名或密码错误
	ErrPasswordWrong = 10101
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = 10102
	// ErrOldPasswordWrong 旧密码错误（修改密码时）
	ErrOldPasswordWrong = 10103

	// --- 用户管理相关错误 ---

	// ErrUsernameDuplicate 用户名已存在
	ErrUsernameDuplicate = 10200
	// ErrEmailDuplicate 邮箱已被注册
	ErrEmailDuplicate = 10201
	// ErrPhoneDuplicate 手机号已被注册
	ErrPhoneDuplicate = 10202

	// --- 角色菜单相关错误 ---

	// ErrRoleCodeDuplicate 角色编码已存在
	ErrRoleCodeDuplicate = 10300
	// ErrRoleNotFound 角色不存在
	ErrRoleNotFound = 10301
	// ErrMenuNotFound 菜单不存在
	ErrMenuNotFound = 10302
)

// ErrCodeMsg 错误码对应的默认消息映射。
var ErrCodeMsg = map[int]string{
	ErrSuccess:           "成功",
	ErrParamInvalid:      "参数校验失败",
	ErrUnauthorized:      "未授权，请先登录",
	ErrForbidden:         "权限不足，无法访问",
	ErrNotFound:          "数据不存在",
	ErrDuplicate:         "数据已存在，不可重复",
	ErrInternal:          "服务内部错误",
	ErrTokenExpired:      "登录已过期，请重新登录",
	ErrTokenInvalid:      "Token无效",
	ErrAccountDisabled:   "账号已被禁用，请联系管理员",
	ErrPasswordWrong:     "用户名或密码错误",
	ErrUserNotFound:      "用户不存在",
	ErrOldPasswordWrong:  "原密码错误",
	ErrUsernameDuplicate: "用户名已存在",
	ErrEmailDuplicate:    "邮箱已被注册",
	ErrPhoneDuplicate:    "手机号已被注册",
	ErrRoleCodeDuplicate: "角色编码已存在",
	ErrRoleNotFound:      "角色不存在",
	ErrMenuNotFound:      "菜单不存在",
}

// CodeError 业务错误类型，实现了 error 接口。
// 携带错误码和错误消息，方便在 logic 层传递具体业务错误。
type CodeError struct {
	// Code 业务错误码
	Code int
	// Msg 错误消息
	Msg string
}

// Error 实现 error 接口，返回格式化的错误字符串。
func (e *CodeError) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", e.Code, e.Msg)
}

// NewCodeError 根据错误码创建 CodeError。
// 消息从预定义的 ErrCodeMsg 映射中查找，找不到则使用"未知错误"。
func NewCodeError(code int) *CodeError {
	msg, ok := ErrCodeMsg[code]
	if !ok {
		msg = "未知错误"
	}
	return &CodeError{Code: code, Msg: msg}
}

// NewCodeErrorMsg 使用自定义消息创建 CodeError。
// 用于错误码相同但需要不同提示信息的场景。
func NewCodeErrorMsg(code int, msg string) *CodeError {
	return &CodeError{Code: code, Msg: msg}
}

// GRPCStatus 实现 gRPC 错误约定。
//
// gRPC 在序列化错误时若发现 error 实现了 GRPCStatus() *status.Status 方法，
// 会直接采用其返回的 status，而不是降级为 codes.Unknown。
// 因此 logic 层只需 return xerr.NewCodeError(...)，业务错误码会无侵入透传到客户端，无需任何手动包装或服务端拦截器。
//
// 业务码同时承载于两处：
//   - status.Code()：映射到最贴近的标准 gRPC 状态码，用于通用中间件分类
//   - status.Details() 的 ErrorInfo.Metadata["code"]：精确还原原始业务码
func (e *CodeError) GRPCStatus() *status.Status {
	st := status.New(toGrpcCode(e.Code), e.Msg)
	stWithDetails, err := st.WithDetails(&errdetails.ErrorInfo{
		Reason: BusinessErrorReason,
		Domain: BusinessErrorDomain,
		Metadata: map[string]string{
			"code": strconv.Itoa(e.Code),
		},
	})
	if err != nil {
		// WithDetails 仅在 proto 序列化失败时报错，对固定字段几乎不可能发生
		return st
	}
	return stWithDetails
}

// toGrpcCode 将业务错误码映射为最贴近的 gRPC 标准状态码。
//
// 该映射只为标准 gRPC 中间件（重试、监控、链路追踪）提供分类依据，
// 真正的业务码通过 ErrorInfo.Metadata["code"] 精确传递，不会因映射收敛而丢失。
func toGrpcCode(code int) codes.Code {
	switch code {
	case ErrParamInvalid:
		return codes.InvalidArgument
	case ErrUnauthorized, ErrTokenExpired, ErrTokenInvalid:
		return codes.Unauthenticated
	case ErrForbidden:
		return codes.PermissionDenied
	case ErrNotFound, ErrUserNotFound, ErrRoleNotFound, ErrMenuNotFound:
		return codes.NotFound
	case ErrDuplicate, ErrUsernameDuplicate, ErrEmailDuplicate, ErrPhoneDuplicate, ErrRoleCodeDuplicate:
		return codes.AlreadyExists
	default:
		return codes.Internal
	}
}
