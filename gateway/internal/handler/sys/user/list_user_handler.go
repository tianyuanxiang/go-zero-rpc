package user

import (
	"net/http"

	"go-zero-rpc/common/response"
	"go-zero-rpc/gateway/internal/logic/sys/user"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListUserReq
		if err := httpx.Parse(r, &req); err != nil {
			response.FailWithMsg(w, r, err.Error())
			return
		}

		l := user.NewListUserLogic(r.Context(), svcCtx)
		resp, err := l.ListUser(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		response.OkWithData(w, r, resp)
	}
}
