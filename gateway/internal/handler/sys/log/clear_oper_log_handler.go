// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package log

import (
	"net/http"

	"go-zero-rpc/common/response"
	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-rpc/gateway/internal/logic/sys/log"
	"go-zero-rpc/gateway/internal/svc"
)

func ClearOperLogHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := log.NewClearOperLogLogic(r.Context(), svcCtx)
		err := l.ClearOperLog()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			response.OK(w, r)
		}
	}
}
