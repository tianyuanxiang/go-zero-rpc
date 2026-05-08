package user

import (
	"net/http"

	"go-zero-rpc/common/response"
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
			// 全局错误处理已在 gateway.go 中通过 httpx.SetErrorHandlerCtx 注册：
			// 自动解析 RPC 透传的 status 错误并以统一 Response 结构返回
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		response.OkWithData(w, r, resp)
	}
}
