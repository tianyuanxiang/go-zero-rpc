// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dict

import (
	"net/http"
	"strconv"

	"go-zero-rpc/common/response"
	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-rpc/gateway/internal/logic/sys/dict"
	"go-zero-rpc/gateway/internal/svc"

	"github.com/zeromicro/go-zero/rest/pathvar"
)

func GetDictDataByTypeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dictTypeStr := pathvar.Vars(r)["dictType"]
		dictTypeId, err := strconv.ParseInt(dictTypeStr, 10, 64)
		if err != nil || dictTypeId <= 0 {
			response.FailWithMsg(w, r, "字典类型ID格式错误")
			return
		}

		l := dict.NewGetDictDataByTypeLogic(r.Context(), svcCtx)
		resp, err := l.GetDictDataByType(dictTypeId)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		response.OkWithData(w, r, resp)
	}
}
