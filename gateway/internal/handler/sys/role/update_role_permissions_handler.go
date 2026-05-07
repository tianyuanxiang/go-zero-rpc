// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-rpc/gateway/internal/logic/sys/role"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
)

func UpdateRolePermissionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateRolePermissionsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := role.NewUpdateRolePermissionsLogic(r.Context(), svcCtx)
		err := l.UpdateRolePermissions(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
