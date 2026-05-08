// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role

import (
	"net/http"
	"strconv"

	"go-zero-rpc/common/response"
	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-rpc/gateway/internal/logic/sys/role"
	"go-zero-rpc/gateway/internal/svc"

	"github.com/zeromicro/go-zero/rest/pathvar"
)

func DeleteRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := pathvar.Vars(r)["id"]
		roleId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || roleId <= 0 {
			response.FailWithMsg(w, r, "角色ID格式错误")
			return
		}

		l := role.NewDeleteRoleLogic(r.Context(), svcCtx)
		err = l.DeleteRole(roleId)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		response.OK(w, r)
	}
}
