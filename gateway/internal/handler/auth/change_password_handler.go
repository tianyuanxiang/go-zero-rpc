// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"net/http"

	"go-zero-rpc/common/response"
	"go-zero-rpc/common/rpcerr"
	"go-zero-rpc/gateway/internal/logic/auth"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ChangePasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChangePasswordReq
		if err := httpx.Parse(r, &req); err != nil {
			response.FailWithMsg(w, r, err.Error())
			return
		}

		if req.OldPassword == "" || req.NewPassword == "" {
			response.FailWithMsg(w, r, "旧密码和新密码不能为空")
			return
		}

		l := auth.NewChangePasswordLogic(r.Context(), svcCtx)
		err := l.ChangePassword(&req)
		if err != nil {
			code, msg := rpcerr.FromStatus(err)
			response.Fail(w, r, code, msg)
			return
		}
		response.OkWithData(w, r, "密码修改成功")
	}
}
