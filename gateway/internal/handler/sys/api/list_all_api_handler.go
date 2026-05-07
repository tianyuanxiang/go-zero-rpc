// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package api

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-rpc/gateway/internal/logic/sys/api"
	"go-zero-rpc/gateway/internal/svc"
)

func ListAllApiHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := api.NewListAllApiLogic(r.Context(), svcCtx)
		resp, err := l.ListAllApi()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
