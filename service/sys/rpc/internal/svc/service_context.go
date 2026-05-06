package svc

import (
	"context"
	"fmt"
	"time"

	"go-zero-rpc/sys-rpc/internal/config"
	sysmodel "go-zero-rpc/sys-rpc/internal/model"
	"go-zero-rpc/sys-rpc/pkg/casbin"
	"go-zero-rpc/sys-rpc/pkg/orm"
	pkgsqlx "go-zero-rpc/sys-rpc/pkg/sqlx"

	casbinv2 "github.com/casbin/casbin/v2"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

// ServiceContext sys-rpc 服务上下文
//
// 与单体版相比的关键变化：
// 1. 不再有 AuthMiddleware/CasbinMiddleware（HTTP 中间件，留在 gateway）
// 2. 不再有业务 Model（如 PlateXXXModel），由各业务 RPC 自己管理
// 3. Enforcer 仍保留：sys-rpc 暴露 CheckPermission 方法供 gateway 远程调用

type ServiceContext struct {
	Config   config.Config
	Orm      *gorm.DB
	RDB      *redis.Client
	Enforcer *casbinv2.Enforcer

	SysUserModel     sysmodel.SysUserModel
	SysRoleModel     sysmodel.SysRoleModel
	SysUserRoleModel sysmodel.SysUserRoleModel
	SysMenuModel     sysmodel.SysMenuModel
	SysRoleMenuModel sysmodel.SysRoleMenuModel
	SysApiModel      sysmodel.SysApiModel
	SysRoleApiModel  sysmodel.SysRoleApiModel
	SysDictTypeModel sysmodel.SysDictTypeModel
	SysDictDataModel sysmodel.SysDictDataModel
	SysLoginLogModel sysmodel.SysLoginLogModel
	SysOperLogModel  sysmodel.SysOperLogModel
	SysFileModel     sysmodel.SysFileModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	rawConn := sqlx.NewMysql(c.DB.DataSource)
	conn := pkgsqlx.NewTimeoutConn(rawConn, 30*time.Second)

	db := orm.NewMysql(&orm.Config{
		DSN:         c.DB.DataSource,
		Active:      20,
		Idle:        10,
		IdleTimeout: time.Hour * 24,
	})

	rdb := redis.NewClient(&redis.Options{
		Addr:        c.BizRedis.Host,
		Password:    c.BizRedis.Pass,
		DB:          c.BizRedis.DB,
		DialTimeout: 5 * time.Second,
		ReadTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Redis 连接失败：%v", err))
	}
	logx.Info("Redis 连接初始化成功")

	enforcer, err := casbin.NewCasbin(db, c.CasbinModelPath)
	if err != nil {
		panic(fmt.Sprintf("Casbin 初始化失败：%v", err))
	}

	return &ServiceContext{
		Config:   c,
		Orm:      db,
		RDB:      rdb,
		Enforcer: enforcer,

		SysUserModel:     sysmodel.NewSysUserModel(conn, c.CacheRedis, db),
		SysRoleModel:     sysmodel.NewSysRoleModel(conn, c.CacheRedis, db),
		SysUserRoleModel: sysmodel.NewSysUserRoleModel(conn, c.CacheRedis, db),
		SysMenuModel:     sysmodel.NewSysMenuModel(conn, c.CacheRedis, db),
		SysRoleMenuModel: sysmodel.NewSysRoleMenuModel(conn, c.CacheRedis, db),
		SysApiModel:      sysmodel.NewSysApiModel(conn, c.CacheRedis, db),
		SysRoleApiModel:  sysmodel.NewSysRoleApiModel(conn, c.CacheRedis, db),
		SysDictTypeModel: sysmodel.NewSysDictTypeModel(conn, c.CacheRedis, db),
		SysDictDataModel: sysmodel.NewSysDictDataModel(conn, c.CacheRedis, db),
		SysLoginLogModel: sysmodel.NewSysLoginLogModel(conn, c.CacheRedis, db),
		SysOperLogModel:  sysmodel.NewSysOperLogModel(conn, c.CacheRedis, db),
		SysFileModel:     sysmodel.NewSysFileModel(conn, c.CacheRedis, db),
	}
}
