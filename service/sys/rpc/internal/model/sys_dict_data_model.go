package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ SysDictDataModel = (*customSysDictDataModel)(nil)

type (
	// SysDictDataModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysDictDataModel.
	SysDictDataModel interface {
		sysDictDataModel
		ListByDictTypeId(ctx context.Context, dictTypeId int64) ([]*SysDictData, int64, error)
		UpdateDictData(ctx context.Context, dictTypeId int64, updates map[string]interface{}) error
		SoftDeleteDictDataByDictTypeIdTrans(ctx context.Context, tx *gorm.DB, dictTypeId int64) error
		SoftDeleteDictDataByDictDataId(ctx context.Context, dictDataId int64) error
	}

	customSysDictDataModel struct {
		*defaultSysDictDataModel
		db *gorm.DB
	}
)

// NewSysDictDataModel returns a model for the database table.
func NewSysDictDataModel(conn sqlx.SqlConn, c cache.CacheConf, db *gorm.DB, opts ...cache.Option) SysDictDataModel {
	return &customSysDictDataModel{
		defaultSysDictDataModel: newSysDictDataModel(conn),
		db:                      db,
	}
}

func (m *customSysDictDataModel) ListByDictTypeId(ctx context.Context, dictTypeId int64) ([]*SysDictData, int64, error) {
	var (
		dictData []*SysDictData
		count    int64
	)
	db := m.db.WithContext(ctx).Table("sys_dict_data").
		Where("type_id = ?", dictTypeId).
		Where("deleted_at is NULL")
	// 查数量
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	// 查数据
	result := db.Find(&dictData)
	if result.Error != nil {
		return nil, count, result.Error
	}
	return dictData, count, nil

}

func (m *customSysDictDataModel) SoftDeleteDictDataByDictTypeIdTrans(ctx context.Context, tx *gorm.DB, dictTypeId int64) error {
	return tx.WithContext(ctx).Table("sys_dict_data").
		Where("type_id = ?", dictTypeId).
		Update("deleted_at", sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}).Error
}

func (m *customSysDictDataModel) SoftDeleteDictDataByDictDataId(ctx context.Context, dictDataId int64) error {
	return m.db.WithContext(ctx).Table("sys_dict_data").
		Where("id = ?", dictDataId).
		Update("deleted_at", sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}).Error
}

func (m *customSysDictDataModel) UpdateDictData(ctx context.Context, dictDataId int64, updates map[string]interface{}) error {
	result := m.db.WithContext(ctx).Table("sys_dict_data").
		Where("id = ? AND deleted_at IS NULL", dictDataId).
		Updates(updates)
	return result.Error
}
