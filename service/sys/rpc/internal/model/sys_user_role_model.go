package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ SysUserRoleModel = (*customSysUserRoleModel)(nil)

type (
	SysUserRoleModel interface {
		sysUserRoleModel
		GetRoleIdsByUserId(ctx context.Context, userId int64) ([]int64, error)
		GetRoleIdsByUserIds(ctx context.Context, userIds []int64) ([]SysUserRole, error)
		AssignRolesTrans(ctx context.Context, tx *gorm.DB, userId int64, roleIds []int64) error
		DeleteUserRoleTrans(ctx context.Context, tx *gorm.DB, roleId int64) error
		DeleteByUserIdTrans(ctx context.Context, tx *gorm.DB, userId int64) error
	}

	customSysUserRoleModel struct {
		*defaultSysUserRoleModel
		db *gorm.DB
	}
)

func NewSysUserRoleModel(conn sqlx.SqlConn, c cache.CacheConf, db *gorm.DB, opts ...cache.Option) SysUserRoleModel {
	return &customSysUserRoleModel{
		defaultSysUserRoleModel: newSysUserRoleModel(conn),
		db:                      db,
	}
}

// GetRoleIdsByUserId 查询指定用户拥有的所有角色ID列表。
func (m *customSysUserRoleModel) GetRoleIdsByUserId(ctx context.Context, userId int64) ([]int64, error) {
	var roleIds []int64
	result := m.db.WithContext(ctx).Table(m.table).
		Where("user_id = ?", userId).
		Pluck("role_id", &roleIds)
	if result.Error != nil {
		return nil, result.Error
	}
	return roleIds, nil
}

// GetRoleIdsByUserIds 查询指定用户拥有的所有角色信息列表
func (m *customSysUserRoleModel) GetRoleIdsByUserIds(ctx context.Context, userIds []int64) ([]SysUserRole, error) {
	var userRoles []SysUserRole
	result := m.db.WithContext(ctx).Table(m.table).
		Where("user_id IN (?)", userIds).
		Find(&userRoles)
	if result.Error != nil {
		return nil, result.Error
	}
	return userRoles, nil
}

func (m *customSysUserRoleModel) AssignRolesTrans(ctx context.Context, tx *gorm.DB, userId int64, roleIds []int64) error {
	// 先删除该用户所有旧的角色关联
	result := tx.WithContext(ctx).Table("sys_user_role").Where("user_id = ?", userId).Delete(&SysUserRole{})
	if result.Error != nil {
		return result.Error
	}

	// 如果没有新角色需要关联，直接返回
	if len(roleIds) == 0 {
		return nil
	}

	// 批量构建新的角色关联
	userRoles := make([]SysUserRole, 0, len(roleIds))
	for _, roleId := range roleIds {
		userRoles = append(userRoles, SysUserRole{
			UserId: userId,
			RoleId: roleId,
		})
	}
	// 批量插入
	return tx.WithContext(ctx).Table("sys_user_role").Create(&userRoles).Error
}

func (m *customSysUserRoleModel) DeleteUserRoleTrans(ctx context.Context, tx *gorm.DB, roleId int64) error {
	return tx.WithContext(ctx).Table("sys_user_role").Where("role_id = ?", roleId).Delete(&SysUserRole{}).Error
}

func (m *customSysUserRoleModel) DeleteByUserIdTrans(ctx context.Context, tx *gorm.DB, userId int64) error {
	return tx.WithContext(ctx).Table("sys_user_role").Where("user_id = ?", userId).Delete(&SysUserRole{}).Error
}
