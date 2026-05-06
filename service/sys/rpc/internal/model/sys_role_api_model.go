package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

var _ SysRoleApiModel = (*customSysRoleApiModel)(nil)

type (
	// SysRoleApiModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysRoleApiModel.
	SysRoleApiModel interface {
		sysRoleApiModel
		InsertRoleApiTrans(ctx context.Context, tx *gorm.DB, roleApi []SysRoleApi) (int64, error)
		DeleteRoleApiByRoleIdTrans(ctx context.Context, tx *gorm.DB, roleId int64) error
		DeleteRoleApiByApiIdTrans(ctx context.Context, tx *gorm.DB, apiId int64) error
		ListRoleIdsByApiId(ctx context.Context, apiId int64) ([]int64, error)
		ListApiIdsByRoleId(ctx context.Context, roleId int64) ([]int64, error)
	}

	customSysRoleApiModel struct {
		*defaultSysRoleApiModel
		db *gorm.DB
	}
)

// NewSysRoleApiModel returns a model for the database table.
func NewSysRoleApiModel(conn sqlx.SqlConn, c cache.CacheConf, db *gorm.DB, opts ...cache.Option) SysRoleApiModel {
	return &customSysRoleApiModel{
		defaultSysRoleApiModel: newSysRoleApiModel(conn),
		db:                     db,
	}
}

func (m *customSysRoleApiModel) InsertRoleApiTrans(ctx context.Context, tx *gorm.DB, roleApi []SysRoleApi) (int64, error) {
	result := tx.WithContext(ctx).Table("sys_role_api").Create(&roleApi)
	return result.RowsAffected, result.Error
}

// DeleteRoleApiByRoleIdTrans 按 role_id 清空指定角色的所有 API 绑定。
// 用于"重新分配角色权限"场景，先清后插，实现全量替换语义
func (m *customSysRoleApiModel) DeleteRoleApiByRoleIdTrans(ctx context.Context, tx *gorm.DB, roleId int64) error {
	return tx.WithContext(ctx).Table("sys_role_api").Where("role_id = ?", roleId).Delete(&SysRoleApi{}).Error
}

// DeleteRoleApiByApiIdTrans 按 api_id 清空指定接口的所有角色绑定。
// 用于"删除 API 接口"场景的级联清理，避免悬挂引用。
// 注意：本方法不用于角色权限重置，错用会导致角色旧绑定未清空触发 uk_role_api 冲突。
func (m *customSysRoleApiModel) DeleteRoleApiByApiIdTrans(ctx context.Context, tx *gorm.DB, apiId int64) error {
	return tx.WithContext(ctx).Table("sys_role_api").Where("api_id = ?", apiId).Delete(&SysRoleApi{}).Error
}

// ListRoleIdsByApiId 根据接口ID查询所有关联的角色ID。
func (m *customSysRoleApiModel) ListRoleIdsByApiId(ctx context.Context, apiId int64) ([]int64, error) {
	var roleIds []int64
	result := m.db.WithContext(ctx).Table("sys_role_api").
		Where("api_id = ?", apiId).
		Pluck("role_id", &roleIds)
	if result.Error != nil {
		return nil, result.Error
	}
	return roleIds, nil
}

// ListApiIdsByRoleId 根据角色ID查询所有关联的接口ID。
func (m *customSysRoleApiModel) ListApiIdsByRoleId(ctx context.Context, roleId int64) ([]int64, error) {
	var apiIds []int64
	result := m.db.WithContext(ctx).Table("sys_role_api").
		Where("role_id = ?", roleId).
		Pluck("api_id", &apiIds)
	if result.Error != nil {
		return nil, result.Error
	}
	return apiIds, nil
}
