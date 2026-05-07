// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-rpc/gateway/internal/logic/sys/role"
	"go-zero-rpc/gateway/internal/svc"
)

func DeleteRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := role.NewDeleteRoleLogic(r.Context(), svcCtx)
		err := l.DeleteRole()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
