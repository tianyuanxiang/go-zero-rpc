package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	// Auth JWT认证配置
	JwtAuth struct {
		// AccessSecret 用于签发和验证访问令牌的密钥，生产环境必须使用强随机字符串
		AccessSecret string
		// AccessExpire 访问令牌过期时间（秒），默认7200（2小时）
		AccessExpire int64
		// RefreshExpire 刷新令牌过期时间（秒），默认604800（7天）
		RefreshExpire int64
	}

	// DB PostgreSQL数据库配置
	DB struct {
		// DataSource PostgreSQL连接字符串，格式：postgres://user:pass@host:port/dbname?sslmode=disable&TimeZone=Asia/Shanghai
		DataSource string
	}

	CacheRedis cache.CacheConf

	// Redis Redis缓存配置
	BizRedis struct {
		// Host Redis地址，格式：host:port
		Host string
		// Pass Redis密码，无密码则留空
		Pass string
		// DB Redis数据库编号，默认0
		DB int
	}

	// CasbinModelPath Casbin RBAC模型配置文件路径，如：etc/rbac_model.conf
	CasbinModelPath string

	// UploadPath 文件上传根目录，如：uploads
	UploadPath string
}
