package model

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ SysFileModel = (*customSysFileModel)(nil)

type (
	// SysFileModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysFileModel.
	SysFileModel interface {
		sysFileModel
		List(ctx context.Context, page, pageSize int, keyword string) ([]*SysFile, int64, error)
		FindActiveById(ctx context.Context, id int64) (*SysFile, error)
		InsertFileReturningId(ctx context.Context, data *SysFile) (int64, error)
		SoftDeleteFile(ctx context.Context, id int64) error
	}

	customSysFileModel struct {
		*defaultSysFileModel
		db *gorm.DB
	}
)

// NewSysFileModel returns a model for the database table.
func NewSysFileModel(conn sqlx.SqlConn, c cache.CacheConf, db *gorm.DB, opts ...cache.Option) SysFileModel {
	return &customSysFileModel{
		defaultSysFileModel: newSysFileModel(conn),
		db:                  db,
	}
}

func (m *customSysFileModel) List(ctx context.Context, page, pageSize int, keyword string) ([]*SysFile, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	db := m.db.WithContext(ctx).Table("sys_file").Where("deleted_at IS NULL")
	if keyword != "" {
		like := "%" + keyword + "%"
		db = db.Where("origin_name ILIKE ? OR filename ILIKE ?", like, like)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	files := make([]*SysFile, 0)
	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

func (m *customSysFileModel) FindActiveById(ctx context.Context, id int64) (*SysFile, error) {
	var file SysFile
	result := m.db.WithContext(ctx).Table("sys_file").
		Where("id = ? AND deleted_at IS NULL", id).
		Take(&file)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, sqlx.ErrNotFound
		}
		return nil, result.Error
	}

	return &file, nil
}

func (m *customSysFileModel) InsertFileReturningId(ctx context.Context, data *SysFile) (int64, error) {
	result := m.db.WithContext(ctx).Table("sys_file").Create(data)
	return data.Id, result.Error
}

func (m *customSysFileModel) SoftDeleteFile(ctx context.Context, id int64) error {
	result := m.db.WithContext(ctx).Table("sys_file").
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return sqlx.ErrNotFound
	}

	return nil
}
