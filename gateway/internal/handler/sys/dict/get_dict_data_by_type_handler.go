// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dict

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-rpc/gateway/internal/logic/sys/dict"
	"go-zero-rpc/gateway/internal/svc"
)

func GetDictDataByTypeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := dict.NewGetDictDataByTypeLogic(r.Context(), svcCtx)
		resp, err := l.GetDictDataByType()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
