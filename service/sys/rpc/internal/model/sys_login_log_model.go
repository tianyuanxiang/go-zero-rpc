package model

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ SysLoginLogModel = (*customSysLoginLogModel)(nil)

type (
	// SysLoginLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysLoginLogModel.
	SysLoginLogModel interface {
		sysLoginLogModel
		List(ctx context.Context, loginLogListReq *LoginLogListReq) ([]SysLoginLog, int64, error)
	}

	customSysLoginLogModel struct {
		*defaultSysLoginLogModel
		db *gorm.DB
	}

	LoginLogListReq struct {
		Page      int
		PageSize  int
		Keyword   string
		Status    int // Status 按登录状态过滤：-1=不过滤，1=成功，0=失败
		StartTime time.Time
		EndTime   time.Time
	}
)

// NewSysLoginLogModel returns a model for the database table.
func NewSysLoginLogModel(conn sqlx.SqlConn, c cache.CacheConf, db *gorm.DB, opts ...cache.Option) SysLoginLogModel {
	return &customSysLoginLogModel{
		defaultSysLoginLogModel: newSysLoginLogModel(conn),
		db:                      db,
	}
}

func (m *customSysLoginLogModel) List(ctx context.Context, req *LoginLogListReq) ([]SysLoginLog, int64, error) {
	var (
		loginLogs []SysLoginLog
		count     int64
	)
	db := m.db.WithContext(ctx).Table("sys_login_log").Where("deleted_at IS NULL")

	// 关键词模糊检索
	if req.Keyword != "" {
		db = db.Where("username like ? OR ip like ? OR location like ? OR browser like ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 时间检索
	if !req.StartTime.IsZero() {
		db = db.Where("login_time >= ?", req.StartTime)
	}
	if !req.EndTime.IsZero() {
		db = db.Where("login_time <= ?", req.EndTime)
	}

	// 状态检索
	if req.Status != -1 {
		db = db.Where("status = ?", req.Status)
	}
	// 查询总数
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := db.Limit(req.PageSize).Offset(offset).Find(&loginLogs).Error; err != nil {
		return nil, 0, err
	}
	return loginLogs, count, nil
}
