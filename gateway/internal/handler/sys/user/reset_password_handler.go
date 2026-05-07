package user

import (
	"net/http"

	"go-zero-rpc/common/response"
	"go-zero-rpc/common/rpcerr"
	"go-zero-rpc/gateway/internal/logic/sys/user"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ResetPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ResetPasswordReq
		if err := httpx.Parse(r, &req); err != nil {
			response.FailWithMsg(w, r, err.Error())
			return
		}

		if req.NewPassword == "" {
			response.FailWithMsg(w, r, "新密码不能为空")
			return
		}

		l := user.NewResetPasswordLogic(r.Context(), svcCtx)
		err := l.ResetPassword(&req)
		if err != nil {
			code, msg := rpcerr.FromStatus(err)
			response.Fail(w, r, code, msg)
			return
		}
		response.OK(w, r)
	}
}
