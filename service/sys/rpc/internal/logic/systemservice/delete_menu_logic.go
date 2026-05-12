// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"
	"go-zero-rpc/sys-rpc/pb"

	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

type DeleteMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuLogic {
	return &DeleteMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteMenu 软删除菜单。
//
// 业务流程：
//  1. 校验菜单存在性
//  2. 校验是否有子菜单（存在则禁止删除）
//  3. 事务内：软删除菜单、删除角色菜单关联
func (l *DeleteMenuLogic) DeleteMenu(in *pb.DeleteMenuReq) (*pb.Empty, error) {
	menuId := in.MenuId

	_, err := l.svcCtx.SysMenuModel.FindOne(l.ctx, menuId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, xerr.NewCodeError(xerr.ErrMenuNotFound)
		}
		l.Errorf("查询菜单[%d]失败：%v", menuId, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	// 检查是否有子菜单
	hasChildren, err := l.svcCtx.SysMenuModel.HasChildren(l.ctx, menuId)
	if err != nil {
		l.Errorf("检查菜单[%d]的子节点失败：%v", menuId, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}
	if hasChildren {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "该菜单存在子菜单，请先删除子菜单")
	}

	// 开启事务
	err = l.svcCtx.Orm.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		// 1.软删除菜单
		if err := l.svcCtx.SysMenuModel.SoftDeleteTrans(l.ctx, tx, menuId); err != nil {
			l.Errorf("删除菜单[%d]失败：%v", menuId, err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		// 2.删除关联关系
		if err := l.svcCtx.SysRoleMenuModel.DeleteRoleMenuByMenuIdTrans(l.ctx, tx, menuId); err != nil {
			l.Errorf("删除角色菜单关联关系失败：%v", err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		return nil
	})
	if err != nil {
		l.Errorf("创建删除菜单事务执行失败: %v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	return &pb.Empty{}, nil
}
