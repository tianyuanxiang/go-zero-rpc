// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"
	"go-zero-rpc/sys-rpc/pb"

	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/svc"
	casbinpkg "go-zero-rpc/sys-rpc/pkg/casbin"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

type DeleteRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteRole 软删除角色，并级联清理用户、菜单、接口关联以及Casbin策略。
//
// 业务流程：
//  1. 校验角色是否存在
//  2. 事务内：软删除角色、清除用户关联、清除菜单关联、清除接口关联
//  3. 事务成功后清除Casbin策略
func (l *DeleteRoleLogic) DeleteRole(in *pb.DeleteRoleReq) (*pb.Empty, error) {
	roleId := in.RoleId

	// 1. 检查角色是否存在
	result, err := l.svcCtx.SysRoleModel.FindOneByRoleId(l.ctx, roleId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, xerr.NewCodeError(xerr.ErrRoleNotFound)
		}
		l.Errorf("查询角色[%d]失败：%v", roleId, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	// 2. 开启事务
	err = l.svcCtx.Orm.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		// 2.1 软删除角色记录
		if err := l.svcCtx.SysRoleModel.SoftDeleteRoleTrans(l.ctx, tx, roleId); err != nil {
			l.Errorf("软删除角色[%d]失败：%v", roleId, err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}

		// 2.2 清除角色与用户的关联
		if err := l.svcCtx.SysUserRoleModel.DeleteUserRoleTrans(l.ctx, tx, roleId); err != nil {
			l.Errorf("清除角色[%d]用户关联失败：%v", roleId, err)
		}

		// 2.3 清除角色与菜单的关联
		if err := l.svcCtx.SysRoleMenuModel.DeleteRoleMenuByRoleIdTrans(l.ctx, tx, roleId); err != nil {
			l.Errorf("清除角色[%d]菜单关联失败：%v", roleId, err)
		}

		// 2.4 清除角色与接口的关联
		if err := l.svcCtx.SysRoleApiModel.DeleteRoleApiByRoleIdTrans(l.ctx, tx, roleId); err != nil {
			l.Errorf("清除角色[%d]接口关联失败：%v", roleId, err)
		}

		return nil
	})
	if err != nil {
		l.Errorf("删除角色事务执行失败: %v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	// 3. 同步Casbin策略
	if err := casbinpkg.RemoveAllPoliciesForRole(l.svcCtx.Enforcer, result.Code); err != nil {
		l.Errorf("删除Casbin策略失败（角色编码: %s）：%v", result.Code, err)
	}

	return &pb.Empty{}, nil
}
