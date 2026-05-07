// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package file

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-rpc/gateway/internal/logic/sys/file"
	"go-zero-rpc/gateway/internal/svc"
)

func DeleteFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := file.NewDeleteFileLogic(r.Context(), svcCtx)
		err := l.DeleteFile()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
