// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"go-zero-rpc/gateway/internal/config"
	"go-zero-rpc/gateway/internal/middleware"
	authclient "go-zero-rpc/sys-rpc/client/authservice"
	permclient "go-zero-rpc/sys-rpc/client/permissionservice"
	systemclient "go-zero-rpc/sys-rpc/client/systemservice"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	AuthRpc authclient.AuthService
	SysRpc  systemclient.SystemService
	PermRpc permclient.PermissionService

	// RDB Redis客户端（用于Token黑名单、缓存等）
	RDB *redis.Client

	AuthMiddleware    rest.Middleware
	CasbinMiddleware  rest.Middleware
	OperLogMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {

	cli := zrpc.MustNewClient(c.SysRpc) // 只建一次，下面三个 client 共用底层 conn

	return &ServiceContext{
		Config:  c,
		AuthRpc: authclient.NewAuthService(cli),
		SysRpc:  systemclient.NewSystemService(cli),
		PermRpc: permclient.NewPermissionService(cli),

		AuthMiddleware:   middleware.NewAuthMiddleware(permclient.NewPermissionService(cli)).Handle(c),
		CasbinMiddleware: middleware.NewCasbinMiddleware(permclient.NewPermissionService(cli)).Handle,
	}
}
