package user

import (
	"net/http"

	"go-zero-rpc/common/response"
	"go-zero-rpc/gateway/internal/logic/sys/user"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CreateUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateUserReq
		if err := httpx.Parse(r, &req); err != nil {
			response.FailWithMsg(w, r, err.Error())
			return
		}

		if req.Username == "" || req.Password == "" {
			response.FailWithMsg(w, r, "用户名和密码不能为空")
			return
		}

		l := user.NewCreateUserLogic(r.Context(), svcCtx)
		err := l.CreateUser(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		response.OK(w, r)
	}
}
