// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"context"
	"fmt"
	"time"

	"go-zero-rpc/gateway/internal/config"
	"go-zero-rpc/gateway/internal/middleware"
	authclient "go-zero-rpc/sys-rpc/client/authservice"
	systemclient "go-zero-rpc/sys-rpc/client/systemservice"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	AuthRpc authclient.AuthService
	SysRpc  systemclient.SystemService

	// RDB Redis客户端（用于Token黑名单、缓存等）
	RDB *redis.Client

	AuthMiddleware    rest.Middleware
	CasbinMiddleware  rest.Middleware
	OperLogMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {

	cli := zrpc.MustNewClient(c.SysRpc) // 只建一次，下面三个 client 共用底层 conn

	// 初始化 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:        c.BizRedis.Host,
		Password:    c.BizRedis.Pass,
		DB:          c.BizRedis.DB,
		DialTimeout: 5 * time.Second,
		ReadTimeout: 3 * time.Second,
	})

	// 验证Redis连接是否正常
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Redis连接失败：%v", err))
	}
	logx.Info("RedisBiz 连接初始化成功")

	return &ServiceContext{
		Config:  c,
		RDB:     rdb,
		AuthRpc: authclient.NewAuthService(cli),
		SysRpc:  systemclient.NewSystemService(cli),

		AuthMiddleware:   middleware.AuthMiddleware(c, rdb),
		CasbinMiddleware: middleware.NewCasbinMiddleware().Handle,
	}
}
