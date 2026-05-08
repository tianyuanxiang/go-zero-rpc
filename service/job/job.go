package main

import (
	"flag"
	"fmt"

	"go-zero-rpc/job/internal/config"
	"go-zero-rpc/job/internal/scheduler"
	"go-zero-rpc/job/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "etc/job.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)

	s := scheduler.NewScheduler(ctx)
	s.Register()

	logx.Infof("starting job service...")
	fmt.Println("Starting job service...")

	s.Start()
}
