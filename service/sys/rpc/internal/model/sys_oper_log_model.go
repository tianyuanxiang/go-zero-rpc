package model

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ SysOperLogModel = (*customSysOperLogModel)(nil)

type (
	// SysOperLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysOperLogModel.
	SysOperLogModel interface {
		sysOperLogModel
		List(ctx context.Context, req *OperLogListReq) ([]SysOperLog, int64, error)
	}

	customSysOperLogModel struct {
		*defaultSysOperLogModel
		db *gorm.DB
	}

	OperLogListReq struct {
		Page         int
		PageSize     int
		Keyword      string // 作用于操作人、请求方式、操作IP
		BusinessType string
		OperType     string
		Status       int // Status 按登录状态过滤：-1=不过滤，1=成功，0=失败
		StartTime    time.Time
		EndTime      time.Time
	}
)

// NewSysOperLogModel returns a model for the database table.
func NewSysOperLogModel(conn sqlx.SqlConn, c cache.CacheConf, db *gorm.DB, opts ...cache.Option) SysOperLogModel {
	return &customSysOperLogModel{
		defaultSysOperLogModel: newSysOperLogModel(conn),
		db:                     db,
	}
}

func (m *customSysOperLogModel) List(ctx context.Context, req *OperLogListReq) ([]SysOperLog, int64, error) {
	var (
		operLogs []SysOperLog
		count    int64
	)
	db := m.db.WithContext(ctx).Table("sys_oper_log").Where("deleted_at IS NULL")

	// 关键词模糊检索
	if req.Keyword != "" {
		db = db.Where("operator_name like ? OR request_method like ? OR oper_ip like ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if req.BusinessType != "" {
		db = db.Where("business_type = ?", req.BusinessType)
	}
	if req.OperType != "" {
		db = db.Where("oper_type = ?", req.OperType)
	}

	// 时间检索
	if !req.StartTime.IsZero() {
		db = db.Where("oper_time >= ?", req.StartTime)
	}
	if !req.EndTime.IsZero() {
		db = db.Where("oper_time <= ?", req.EndTime)
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
	if err := db.Limit(req.PageSize).Offset(offset).Find(&operLogs).Error; err != nil {
		return nil, 0, err
	}
	return operLogs, count, nil
}
