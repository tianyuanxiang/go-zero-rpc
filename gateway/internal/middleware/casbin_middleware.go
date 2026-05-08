// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package middleware

import (
	"go-zero-rpc/common/response"
	permclient "go-zero-rpc/sys-rpc/client/permissionservice"
	"go-zero-rpc/sys-rpc/sys"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type CasbinMiddleware struct {
	permRpc permclient.PermissionService
}

func NewCasbinMiddleware(permRpc permclient.PermissionService) *CasbinMiddleware {
	return &CasbinMiddleware{permRpc: permRpc}
}

func (m *CasbinMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := GetUserIdFromCtx(r.Context())
		if userId == 0 {
			response.FailUnauthorized(w, r)
			return
		}
		resp, err := m.permRpc.CheckPermission(r.Context(), &sys.CheckPermissionReq{
			UserId: userId,
			Path:   r.URL.Path,
			Method: r.Method,
		})
		if err != nil {
			logx.WithContext(r.Context()).Errorf(
				"CheckPermission RPC failed! user=%d path=%s method=%s err=%v",
				userId, r.URL.Path, r.Method, err,
			)
			response.FailForbidden(w, r)
			return
		}
		if !resp.Allowed {
			logx.WithContext(r.Context()).Infof(
				"Authorization denied! user=%d path=%s method=%s reason=%s",
				userId, r.URL.Path, r.Method, resp.Reason,
			)
			response.FailForbidden(w, r)
			return
		}
		
		next.ServeHTTP(w, r)
	}
}
