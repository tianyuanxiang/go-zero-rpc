package user

import (
	"net/http"

	"go-zero-rpc/common/response"
	"go-zero-rpc/common/rpcerr"
	"go-zero-rpc/gateway/internal/logic/sys/user"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteUserReq
		if err := httpx.Parse(r, &req); err != nil {
			response.FailWithMsg(w, r, err.Error())
			return
		}

		l := user.NewDeleteUserLogic(r.Context(), svcCtx)
		err := l.DeleteUser(&req)
		if err != nil {
			code, msg := rpcerr.FromStatus(err)
			response.Fail(w, r, code, msg)
			return
		}
		response.OK(w, r)
	}
}
