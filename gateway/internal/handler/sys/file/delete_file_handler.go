// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package file

import (
	"net/http"
	"strconv"

	"go-zero-rpc/common/response"
	"github.com/zeromicro/go-zero/rest/httpx"
	"go-zero-rpc/gateway/internal/logic/sys/file"
	"go-zero-rpc/gateway/internal/svc"

	"github.com/zeromicro/go-zero/rest/pathvar"
)

func DeleteFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := pathvar.Vars(r)["id"]
		fileId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || fileId <= 0 {
			response.FailWithMsg(w, r, "文件ID格式错误")
			return
		}

		l := file.NewDeleteFileLogic(r.Context(), svcCtx)
		err = l.DeleteFile(fileId)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		response.OK(w, r)
	}
}
