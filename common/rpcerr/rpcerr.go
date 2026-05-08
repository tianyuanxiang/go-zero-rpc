// Package rpcerr 提供 gRPC 错误与业务错误码之间的客户端侧解析工具。
//
// 服务端无需本包：业务 logic 直接 return xerr.NewCodeError(...) 即可，
// gRPC 框架会通过 *CodeError.GRPCStatus() 自动序列化错误，
// 业务码透传机制完全由 xerr 包实现。
package rpcerr

import (
	"strconv"

	"go-zero-rpc/common/xerr"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FromStatus 从 gRPC 客户端收到的 error 中解析出业务码与消息。
//
// 解析顺序：
//  1. 优先从 status.Details() 的 ErrorInfo.Metadata["code"] 中提取精确业务码；
//     这是本项目 RPC 服务通过 *CodeError.GRPCStatus() 写入的字段，能 100% 还原。
//  2. 若没有 ErrorInfo（来自非本项目 RPC，或服务端未实现 GRPCStatus），
//     则按标准 gRPC 状态码映射到通用业务码。
//  3. 任何分支都保留 status.Message() 作为提示，避免吞掉服务端原始报错。
func FromStatus(err error) (int, string) {
	if err == nil {
		return xerr.ErrSuccess, xerr.ErrCodeMsg[xerr.ErrSuccess]
	}

	st, ok := status.FromError(err)
	if !ok {
		// 非 gRPC 错误（例如本地业务 err、网络层 err），直接透出原文便于排查
		return xerr.ErrInternal, err.Error()
	}

	for _, detail := range st.Details() {
		info, ok := detail.(*errdetails.ErrorInfo)
		if !ok || info.Reason != xerr.BusinessErrorReason {
			continue
		}
		code, convErr := strconv.Atoi(info.Metadata["code"])
		if convErr == nil {
			return code, st.Message()
		}
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return xerr.ErrParamInvalid, st.Message()
	case codes.Unauthenticated:
		return xerr.ErrUnauthorized, st.Message()
	case codes.PermissionDenied:
		return xerr.ErrForbidden, st.Message()
	case codes.NotFound:
		return xerr.ErrNotFound, st.Message()
	case codes.AlreadyExists:
		return xerr.ErrDuplicate, st.Message()
	default:
		// 未识别的 grpc code（含 codes.Unknown / Internal 等）：
		// 业务码归一到 ErrInternal，但 message 保留 status 原文，
		// 方便定位「服务端没实现 GRPCStatus」「上游异常」等问题。
		return xerr.ErrInternal, st.Message()
	}
}
