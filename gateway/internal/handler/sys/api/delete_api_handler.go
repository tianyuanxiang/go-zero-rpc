// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package api

import (
	"net/http"
	"strconv"

	"go-zero-rpc/common/response"
	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-rpc/gateway/internal/logic/sys/api"
	"go-zero-rpc/gateway/internal/svc"

	"github.com/zeromicro/go-zero/rest/pathvar"
)

func DeleteApiHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := pathvar.Vars(r)["id"]
		apiId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || apiId <= 0 {
			response.FailWithMsg(w, r, "API ID格式错误")
			return
		}

		l := api.NewDeleteApiLogic(r.Context(), svcCtx)
		err = l.DeleteApi(apiId)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		response.OK(w, r)
	}
}
