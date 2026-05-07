// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dict

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-rpc/gateway/internal/logic/sys/dict"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
)

func UpdateDictDataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateDictDataReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := dict.NewUpdateDictDataLogic(r.Context(), svcCtx)
		err := l.UpdateDictData(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
