// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package menu

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-rpc/gateway/internal/logic/sys/menu"
	"go-zero-rpc/gateway/internal/svc"
)

func GetCurrentUserMenusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := menu.NewGetCurrentUserMenusLogic(r.Context(), svcCtx)
		resp, err := l.GetCurrentUserMenus()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
