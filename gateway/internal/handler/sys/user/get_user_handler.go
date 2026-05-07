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

func GetUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserReq
		if err := httpx.Parse(r, &req); err != nil {
			response.FailWithMsg(w, r, err.Error())
			return
		}

		l := user.NewGetUserLogic(r.Context(), svcCtx)
		resp, err := l.GetUser(&req)
		if err != nil {
			code, msg := rpcerr.FromStatus(err)
			response.Fail(w, r, code, msg)
			return
		}
		response.OkWithData(w, r, resp)
	}
}
