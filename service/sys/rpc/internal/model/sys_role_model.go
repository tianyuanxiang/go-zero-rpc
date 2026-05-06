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

var _ SysRoleModel = (*customSysRoleModel)(nil)

type (
	SysRoleModel interface {
		sysRoleModel
		List(ctx context.Context, page, pageSize int, keyword string) ([]*SysRole, int64, error)
		ListAll(ctx context.Context) ([]*SysRole, error)
		FindOneByRoleId(ctx context.Context, roleId int64) (*SysRole, error)
		FindByIds(ctx context.Context, roleIds []int64) ([]*SysRole, error)
		InsertRoleTrans(ctx context.Context, tx *gorm.DB, role *SysRole) (int64, error)
		UpdateRoleTrans(ctx context.Context, id int64, updates map[string]interface{}) error
		SoftDeleteRoleTrans(ctx context.Context, tx *gorm.DB, roleId int64) error
	}

	customSysRoleModel struct {
		*defaultSysRoleModel
		db *gorm.DB
	}
)

func NewSysRoleModel(conn sqlx.SqlConn, c cache.CacheConf, db *gorm.DB, opts ...cache.Option) SysRoleModel {
	return &customSysRoleModel{
		defaultSysRoleModel: newSysRoleModel(conn),
		db:                  db,
	}
}

// 查询角色列表
func (m *customSysRoleModel) List(ctx context.Context, page, pageSize int, keyword string) ([]*SysRole, int64, error) {
	db := m.db.WithContext(ctx).Table("sys_role").Where("deleted_at IS NULL")
	// 关键词模糊检索
	if keyword != "" {
		db = db.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	// 查询总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	// 分页查询
	var roles []*SysRole
	offset := (page - 1) * pageSize
	if err := db.Limit(pageSize).Offset(offset).Find(&roles).Error; err != nil {
		return nil, 0, err
	}
	return roles, total, nil
}

// ListAll 查询所有启用状态的角色（不分页，用于下拉选择）。
func (m *customSysRoleModel) ListAll(ctx context.Context) ([]*SysRole, error) {
	var roles []*SysRole
	result := m.db.WithContext(ctx).Table("sys_role").
		Select("id", "name", "code").
		Where("deleted_at IS NULL AND status = 1").
		Order("sort ASC, id ASC").
		Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}

// FindOneByRoleId 按角色ID查询未软删除的角色。
func (m *customSysRoleModel) FindOneByRoleId(ctx context.Context, roleId int64) (*SysRole, error) {
	var role SysRole
	result := m.db.WithContext(ctx).Table("sys_role").
		Where("id = ? AND deleted_at IS NULL", roleId).
		First(&role)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &role, nil
}

// FindByIds 批量查询多个角色信息
func (m *customSysRoleModel) FindByIds(ctx context.Context, roleIds []int64) ([]*SysRole, error) {
	var roles []*SysRole
	result := m.db.WithContext(ctx).Table("sys_role").
		Where("id IN ? AND deleted_at IS NULL", roleIds).
		Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}

func (m *customSysRoleModel) InsertRoleTrans(ctx context.Context, tx *gorm.DB, role *SysRole) (int64, error) {
	result := tx.WithContext(ctx).Table("sys_role").Create(&role)

	return role.Id, result.Error
}

// SoftDeleteRoleTrans 在事务中软删除角色。
func (m *customSysRoleModel) SoftDeleteRoleTrans(ctx context.Context, tx *gorm.DB, roleId int64) error {
	result := tx.WithContext(ctx).Table("sys_role").
		Where("id = ?", roleId).
		Update("deleted_at", sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		})
	return result.Error
}

func (m *customSysRoleModel) UpdateRoleTrans(ctx context.Context, id int64, updates map[string]interface{}) error {
	result := m.db.WithContext(ctx).Table("sys_role").
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates)
	return result.Error
}
