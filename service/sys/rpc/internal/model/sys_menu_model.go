package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ SysMenuModel = (*customSysMenuModel)(nil)

type (
	SysMenuModel interface {
		sysMenuModel
		ListByIds(ctx context.Context, ids []int64) ([]*SysMenu, error)
		ListAll(ctx context.Context) ([]*SysMenu, error)
		// HasChildren 检查是否有子菜单
		HasChildren(ctx context.Context, id int64) (bool, error)
		InsertMenuReturningId(ctx context.Context, data *SysMenu) (int64, error)
		SoftDeleteTrans(ctx context.Context, tx *gorm.DB, id int64) error
		UpdateMenuTrans(ctx context.Context, tx *gorm.DB, id int64, updates map[string]interface{}) error
		BatchUpdateMenuStatus(ctx context.Context, tx *gorm.DB, ids []int64, status int) error
		BatchUpdateMenuVisible(ctx context.Context, tx *gorm.DB, ids []int64, visable int) error
	}

	customSysMenuModel struct {
		*defaultSysMenuModel
		db *gorm.DB
	}
)

func NewSysMenuModel(conn sqlx.SqlConn, c cache.CacheConf, db *gorm.DB, opts ...cache.Option) SysMenuModel {
	return &customSysMenuModel{
		defaultSysMenuModel: newSysMenuModel(conn),
		db:                  db,
	}
}

// ListByIds 根据菜单ID列表批量查询未软删除的菜单，按父级和排序字段升序返回。
func (m *customSysMenuModel) ListByIds(ctx context.Context, ids []int64) ([]*SysMenu, error) {
	if len(ids) == 0 {
		return []*SysMenu{}, nil
	}
	var menus []*SysMenu
	result := m.db.WithContext(ctx).Table("sys_menu").
		Where("id IN ? AND deleted_at IS NULL", ids).
		Order("parent_id ASC, sort ASC").
		Find(&menus)
	if result.Error != nil {
		return nil, result.Error
	}
	return menus, nil
}

func (m *customSysMenuModel) ListAll(ctx context.Context) ([]*SysMenu, error) {
	var menus []*SysMenu
	result := m.db.WithContext(ctx).Table("sys_menu").
		Where("deleted_at IS NULL").
		Order("parent_id ASC, sort ASC").
		Find(&menus)
	if result.Error != nil {
		return nil, result.Error
	}
	return menus, nil
}
func (m *customSysMenuModel) HasChildren(ctx context.Context, id int64) (bool, error) {
	var count int64

	result := m.db.WithContext(ctx).Table("sys_menu").
		Where("parent_id = ?", id).
		Where("deleted_at IS NULL").Count(&count)

	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func (m *customSysMenuModel) InsertMenuReturningId(ctx context.Context, data *SysMenu) (int64, error) {
	result := m.db.WithContext(ctx).Table("sys_menu").Create(data)
	return data.Id, result.Error
}

func (m *customSysMenuModel) SoftDeleteTrans(ctx context.Context, tx *gorm.DB, id int64) error {
	return tx.WithContext(ctx).Table("sys_menu").
		Where("id = ?", id).
		Update("deleted_at", sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}).Error
}

func (m *customSysMenuModel) UpdateMenuTrans(ctx context.Context, tx *gorm.DB, id int64, updates map[string]interface{}) error {
	result := tx.WithContext(ctx).Table("sys_menu").
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates)
	return result.Error
}

func (m *customSysMenuModel) BatchUpdateMenuStatus(ctx context.Context, tx *gorm.DB, ids []int64, status int) error {
	return tx.WithContext(ctx).Table("sys_menu").
		Where("id IN ? AND deleted_at IS NULL", ids).
		Update("status", status).Error
}

func (m *customSysMenuModel) BatchUpdateMenuVisible(ctx context.Context, tx *gorm.DB, ids []int64, visible int) error {
	return tx.WithContext(ctx).Table("sys_menu").
		Where("id IN ? AND deleted_at IS NULL", ids).
		Update("visible", visible).Error
}
