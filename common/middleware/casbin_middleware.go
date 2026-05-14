package middleware

import (
	"net/http"

	"go-zero-rpc/common/response"
	permclient "go-zero-rpc/sys-rpc/client/permissionservice"
	"go-zero-rpc/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

// CasbinMiddleware RBAC权限校验中间件。
type CasbinMiddleware struct {
	permRpc permclient.PermissionService
}

// NewCasbinMiddleware 创建权限校验中间件。
func NewCasbinMiddleware(permRpc permclient.PermissionService) *CasbinMiddleware {
	return &CasbinMiddleware{permRpc: permRpc}
}

// Handle 返回Casbin权限校验中间件函数。
func (m *CasbinMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := GetUserIdFromCtx(r.Context())
		if userId == 0 {
			response.FailUnauthorized(w, r)
			return
		}
		resp, err := m.permRpc.CheckPermission(r.Context(), &pb.CheckPermissionReq{
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
