package model

import (
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
