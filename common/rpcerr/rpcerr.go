package rpcerr

import (
	"strconv"

	"go-zero-rpc/common/xerr"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const businessReason = "BUSINESS_ERROR"

func ToStatus(err error) error {
	if err == nil {
		return nil
	}

	codeErr, ok := err.(*xerr.CodeError)
	if !ok {
		st := status.New(codes.Internal, xerr.ErrCodeMsg[xerr.ErrInternal])
		return st.Err()
	}

	grpcCode := toGrpcCode(codeErr.Code)
	st := status.New(grpcCode, codeErr.Msg)
	stWithDetails, detailErr := st.WithDetails(&errdetails.ErrorInfo{
		Reason: businessReason,
		Domain: "go-zero-rpc",
		Metadata: map[string]string{
			"code": strconv.Itoa(codeErr.Code),
		},
	})
	if detailErr != nil {
		return st.Err()
	}
	return stWithDetails.Err()
}

func FromStatus(err error) (int, string) {
	if err == nil {
		return xerr.ErrSuccess, xerr.ErrCodeMsg[xerr.ErrSuccess]
	}

	st, ok := status.FromError(err)
	if !ok {
		return xerr.ErrInternal, xerr.ErrCodeMsg[xerr.ErrInternal]
	}

	for _, detail := range st.Details() {
		info, ok := detail.(*errdetails.ErrorInfo)
		if !ok || info.Reason != businessReason {
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
	default:
		return xerr.ErrInternal, xerr.ErrCodeMsg[xerr.ErrInternal]
	}
}

func toGrpcCode(code int) codes.Code {
	switch code {
	case xerr.ErrParamInvalid:
		return codes.InvalidArgument
	case xerr.ErrUnauthorized, xerr.ErrTokenExpired, xerr.ErrTokenInvalid:
		return codes.Unauthenticated
	case xerr.ErrForbidden:
		return codes.PermissionDenied
	case xerr.ErrNotFound, xerr.ErrUserNotFound, xerr.ErrRoleNotFound, xerr.ErrMenuNotFound:
		return codes.NotFound
	case xerr.ErrDuplicate, xerr.ErrUsernameDuplicate, xerr.ErrEmailDuplicate, xerr.ErrPhoneDuplicate, xerr.ErrRoleCodeDuplicate:
		return codes.AlreadyExists
	default:
		return codes.Internal
	}
}
