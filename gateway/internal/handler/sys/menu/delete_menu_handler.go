// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package menu

import (
	"net/http"
	"strconv"

	"go-zero-rpc/common/response"
	"go-zero-rpc/gateway/internal/logic/sys/menu"
	"go-zero-rpc/gateway/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

func DeleteMenuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := pathvar.Vars(r)["id"]
		menuId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || menuId <= 0 {
			response.FailWithMsg(w, r, "菜单ID格式错误")
			return
		}

		l := menu.NewDeleteMenuLogic(r.Context(), svcCtx)
		err = l.DeleteMenu(menuId)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		response.OK(w, r)
	}
}
