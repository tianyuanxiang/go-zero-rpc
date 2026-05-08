// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	SysRpc zrpc.RpcClientConf // <-- 关键：通过 etcd 发现 sys-rpc

	Upload struct {
		Path        string
		MaxSize     int64    // 字节
		AllowedExts []string // 例 [".png", ".jpg", ".pdf"]
	}
}
