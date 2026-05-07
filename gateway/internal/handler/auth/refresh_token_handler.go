// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"go-zero-rpc/common/response"
	"go-zero-rpc/common/rpcerr"
	"net/http"

	"go-zero-rpc/gateway/internal/logic/auth"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func RefreshTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RefreshTokenReq
		if err := httpx.Parse(r, &req); err != nil {
			response.FailWithMsg(w, r, err.Error())
			return
		}

		if req.RefreshToken == "" {
			response.FailWithMsg(w, r, "RefreshToken不能为空")
			return
		}

		l := auth.NewRefreshTokenLogic(r.Context(), svcCtx)
		resp, err := l.RefreshToken(&req)
		if err != nil {
			code, msg := rpcerr.FromStatus(err)
			response.Fail(w, r, code, msg)
			return
		} else {
			response.OkWithData(w, r, resp)
		}
	}
}
