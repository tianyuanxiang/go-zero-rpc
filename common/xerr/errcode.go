package xerr

import "fmt"

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

	// --- 业务数据相关错误 ---

	// ErrTankNotFound 槽体不存在
	ErrTankNotFound = 10400
	// ErrTankIDDuplicate 槽体ID已存在
	ErrTankIDDuplicate = 10401
	// ErrNoStateFound 模型状态数据不存在（槽体未初始化）
	ErrNoStateFound = 10402
	// ErrEventNotFound 事件记录不存在
	ErrEventNotFound = 10403
	// ErrFileUploadFailed 文件上传失败
	ErrFileUploadFailed = 10404
	// ErrExcelParseFailed Excel解析失败
	ErrExcelParseFailed = 10405
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
	ErrTankNotFound:      "槽体不存在",
	ErrTankIDDuplicate:   "槽体ID已存在",
	ErrNoStateFound:      "模型状态不存在，请先初始化槽体",
	ErrEventNotFound:     "事件记录不存在",
	ErrFileUploadFailed:  "文件上传失败",
	ErrExcelParseFailed:  "Excel文件解析失败",
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
