package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ SysDictTypeModel = (*customSysDictTypeModel)(nil)

type (
	// SysDictTypeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysDictTypeModel.
	SysDictTypeModel interface {
		sysDictTypeModel
		SoftDeleteDictTypeTrans(ctx context.Context, tx *gorm.DB, dictTypeId int64) error
		List(ctx context.Context, page, pageSize int, keyword string) ([]*SysDictType, int64, error)
		InsertDictTypeReturningId(ctx context.Context, data *SysDictType) (int64, error)
		UpdateDictType(ctx context.Context, id int64, updates map[string]interface{}) error
	}

	customSysDictTypeModel struct {
		*defaultSysDictTypeModel
		db *gorm.DB
	}
)

// NewSysDictTypeModel returns a model for the database table.
func NewSysDictTypeModel(conn sqlx.SqlConn, c cache.CacheConf, db *gorm.DB, opts ...cache.Option) SysDictTypeModel {
	return &customSysDictTypeModel{
		defaultSysDictTypeModel: newSysDictTypeModel(conn),
		db:                      db,
	}
}

func (m *customSysDictTypeModel) List(ctx context.Context, page, pageSize int, keyword string) ([]*SysDictType, int64, error) {
	db := m.db.WithContext(ctx).Table("sys_dict_type").Where("deleted_at IS NULL")
	// 关键词模糊检索
	if keyword != "" {
		db = db.Where("name like ? OR code like ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	// 查询总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	// 分页查询
	var dictTypes []*SysDictType
	offset := (page - 1) * pageSize
	if err := db.Limit(pageSize).Offset(offset).Find(&dictTypes).Error; err != nil {
		return nil, 0, err
	}
	return dictTypes, total, nil
}

func (m *customSysDictTypeModel) InsertDictTypeReturningId(ctx context.Context, data *SysDictType) (int64, error) {
	result := m.db.WithContext(ctx).Table("sys_dict_type").Create(data)
	return data.Id, result.Error
}

func (m *customSysDictTypeModel) SoftDeleteDictTypeTrans(ctx context.Context, tx *gorm.DB, dictTypeId int64) error {
	return tx.WithContext(ctx).Table("sys_dict_type").
		Where("id = ?", dictTypeId).
		Update("deleted_at", sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}).Error
}

func (m *customSysDictTypeModel) UpdateDictType(ctx context.Context, dictTypeId int64, updates map[string]interface{}) error {
	result := m.db.WithContext(ctx).Table("sys_dict_type").
		Where("id = ? AND deleted_at IS NULL", dictTypeId).
		Updates(updates)
	return result.Error
}
