// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"go-zero-rpc/common/response"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"

	"go-zero-rpc/gateway/internal/logic/auth"
	"go-zero-rpc/gateway/internal/svc"
)

func LogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := auth.NewLogoutLogic(r.Context(), svcCtx)
		err := l.Logout(r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		} else {
			response.OK(w, r)
		}
	}
}
