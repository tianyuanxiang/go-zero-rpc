package main

import (
	"flag"
	"fmt"
	"time"

	"go-zero-rpc/sys-rpc/internal/config"

	authsvr "go-zero-rpc/sys-rpc/internal/server/authservice"
	"go-zero-rpc/sys-rpc/internal/svc"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/sys.yaml", "the config file")

func main() {
	flag.Parse()

	logx.DisableStat()
	// 将慢查询阈值设为 1ms，使所有 SQL 语句都被记录到日志中（仅开发调试用）
	sqlx.SetSlowThreshold(time.Millisecond)

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		sys.RegisterAuthServiceServer(grpcServer, authsvr.NewAuthServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
