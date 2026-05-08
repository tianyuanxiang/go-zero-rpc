// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dict

import (
	"net/http"
	"strconv"

	"go-zero-rpc/common/response"
	"go-zero-rpc/gateway/internal/logic/sys/dict"
	"go-zero-rpc/gateway/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

func DeleteDictTypeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := pathvar.Vars(r)["id"]
		dictTypeId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || dictTypeId <= 0 {
			response.FailWithMsg(w, r, "字典类型ID格式错误")
			return
		}

		l := dict.NewDeleteDictTypeLogic(r.Context(), svcCtx)
		err = l.DeleteDictType(dictTypeId)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		response.OK(w, r)
	}
}
