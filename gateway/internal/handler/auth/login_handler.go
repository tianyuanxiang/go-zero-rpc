// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"go-zero-rpc/common/response"
	"go-zero-rpc/gateway/internal/logic/auth"
	"net/http"

	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			response.FailWithMsg(w, r, err.Error())
			return
		}

		if req.Username == "" || req.Password == "" {
			response.FailWithMsg(w, r, "用户名和密码不能为空")
			return
		}

		l := auth.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req, r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		response.OkWithData(w, r, resp)

	}
}
