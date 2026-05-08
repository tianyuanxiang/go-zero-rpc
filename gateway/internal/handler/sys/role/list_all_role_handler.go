// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-rpc/common/response"
	"go-zero-rpc/gateway/internal/logic/sys/role"
	"go-zero-rpc/gateway/internal/svc"
)

func ListAllRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := role.NewListAllRoleLogic(r.Context(), svcCtx)
		resp, err := l.ListAllRole()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			response.OkWithData(w, r, resp)
		}
	}
}
