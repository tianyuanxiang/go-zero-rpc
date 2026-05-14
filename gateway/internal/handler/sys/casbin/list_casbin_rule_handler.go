// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package casbin

import (
	"go-zero-rpc/common/response"
	"net/http"

	"go-zero-rpc/gateway/internal/logic/sys/casbin"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListCasbinRuleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListCasbinRuleReq
		if err := httpx.Parse(r, &req); err != nil {
			response.FailWithMsg(w, r, err.Error())
			return
		}

		l := casbin.NewListCasbinRuleLogic(r.Context(), svcCtx)
		resp, err := l.ListCasbinRule(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			response.OkWithData(w, r, resp)
		}
	}
}
