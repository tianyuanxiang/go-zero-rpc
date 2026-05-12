// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	commonmw "go-zero-rpc/common/middleware"
	"go-zero-rpc/gateway/internal/config"
	"go-zero-rpc/gateway/internal/ws"
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

	// WsHub WebSocket 连接管理器
	WsHub *ws.Hub

	AuthMiddleware    rest.Middleware
	CasbinMiddleware  rest.Middleware
	OperLogMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {

	cli := zrpc.MustNewClient(c.SysRpc) // 只建一次，下面三个 client 共用底层 conn

	// 创建 WebSocket Hub
	wsHub := ws.NewHub()
	go wsHub.Run()

	return &ServiceContext{
		Config:  c,
		AuthRpc: authclient.NewAuthService(cli),
		SysRpc:  systemclient.NewSystemService(cli),
		PermRpc: permclient.NewPermissionService(cli),
		WsHub:   wsHub,

		AuthMiddleware:   commonmw.NewAuthMiddleware(permclient.NewPermissionService(cli)).Handle(c.Auth.AccessSecret),
		CasbinMiddleware: commonmw.NewCasbinMiddleware(permclient.NewPermissionService(cli)).Handle,
	}
}
