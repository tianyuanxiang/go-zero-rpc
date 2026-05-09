package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ SysApiModel = (*customSysApiModel)(nil)

type (
	// SysApiModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysApiModel.
	SysApiModel interface {
		sysApiModel
		ListByIds(ctx context.Context, ids []int64) ([]SysApi, error)
		List(ctx context.Context, page, pageSize int, group, keyword string) ([]*SysApi, int64, error)
		InsertApiReturningId(ctx context.Context, data *SysApi) (int64, error)
		UpdateApi(ctx context.Context, id int64, updates map[string]interface{}) error
		SoftDeleteApiTrans(ctx context.Context, tx *gorm.DB, apiIds int64) error
	}

	customSysApiModel struct {
		*defaultSysApiModel
		db *gorm.DB
	}
)

// NewSysApiModel returns a model for the database table.
func NewSysApiModel(conn sqlx.SqlConn, c cache.CacheConf, db *gorm.DB, opts ...cache.Option) SysApiModel {
	return &customSysApiModel{
		defaultSysApiModel: newSysApiModel(conn),
		db:                 db,
	}
}

// 查询api列表
func (m *customSysApiModel) List(ctx context.Context, page, pageSize int, group, keyword string) ([]*SysApi, int64, error) {
	db := m.db.WithContext(ctx).Table("sys_api").Where("deleted_at IS NULL")
	// 关键词模糊检索
	if keyword != "" {
		db = db.Where("api_path like ? OR method like ? OR description like ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if group != "" {
		db = db.Where("api_group like ?", "%"+group+"%")
	}

	// 查询总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	// 分页查询
	var apis []*SysApi
	offset := (page - 1) * pageSize
	if err := db.Limit(pageSize).Offset(offset).Find(&apis).Error; err != nil {
		return nil, 0, err
	}
	return apis, total, nil

}

func (m *customSysApiModel) ListByIds(ctx context.Context, ids []int64) ([]SysApi, error) {
	var apis []SysApi

	db := m.db.WithContext(ctx).Table("sys_api").Where("deleted_at IS NULL")
	if len(ids) > 0 {
		db = db.Where("id in ?", ids)
	}
	result := db.Find(&apis)
	if result.Error != nil {
		return nil, result.Error
	}
	return apis, nil
}

func (m *customSysApiModel) InsertApiReturningId(ctx context.Context, data *SysApi) (int64, error) {
	result := m.db.WithContext(ctx).Table("sys_api").Create(data)
	return data.Id, result.Error
}

func (m *customSysApiModel) SoftDeleteApiTrans(ctx context.Context, tx *gorm.DB, apiIds int64) error {
	return tx.WithContext(ctx).Table("sys_api").
		Where("id = ?", apiIds).
		Update("deleted_at", sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}).Error
}

func (m *customSysApiModel) UpdateApi(ctx context.Context, id int64, updates map[string]interface{}) error {
	result := m.db.WithContext(ctx).Table("sys_api").
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates)
	return result.Error
}
