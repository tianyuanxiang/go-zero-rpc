// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"

	"go-zero-rpc/common/xerr"
	systemmodel "go-zero-rpc/sys-rpc/internal/model"
	"go-zero-rpc/sys-rpc/internal/svc"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

type UpdateMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuLogic {
	return &UpdateMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateMenu 更新菜单基本信息，状态/可见性变化时级联子孙节点。
//
// 业务流程：
//  1. 查询菜单存在性
//  2. 拼接动态更新字段
//  3. 事务内：状态/可见性变化时级联更新子孙；执行菜单更新
func (l *UpdateMenuLogic) UpdateMenu(in *sys.UpdateMenuReq) (*sys.Empty, error) {
	oldMenu, err := l.svcCtx.SysMenuModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, xerr.NewCodeError(xerr.ErrMenuNotFound)
		}
		l.Errorf("查询菜单[%d]失败：%v", in.Id, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}
	if oldMenu.DeletedAt.Valid {
		return nil, xerr.NewCodeError(xerr.ErrMenuNotFound)
	}

	// 拼接动态更新字段
	updates := make(map[string]interface{})
	if in.HasParentId {
		updates["parent_id"] = in.ParentId
	}
	if in.HasMenuName {
		updates["name"] = in.MenuName
	}
	if in.HasPath {
		updates["menu_path"] = in.Path
	}
	if in.HasComponent {
		updates["component"] = in.Component
	}
	if in.HasIcon {
		updates["icon"] = in.Icon
	}
	if in.HasRemark {
		updates["remark"] = in.Remark
	}
	if in.HasMenuType {
		updates["menu_type"] = in.MenuType
	}
	if in.HasSort {
		updates["sort"] = in.Sort
	}
	if in.HasPerms {
		updates["permission"] = in.Perms
	}
	if in.HasVisible {
		updates["visible"] = in.Visible
	}
	if in.HasStatus {
		updates["status"] = in.Status
	}

	if len(updates) == 0 {
		l.Error("更新字段为空")
		return nil, xerr.NewCodeError(xerr.ErrParamInvalid)
	}

	// 开启事务
	err = l.svcCtx.Orm.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		// 1.1 停用级联子孙
		if in.HasStatus && in.Status == 0 && oldMenu.Status == 1 {
			if err := l.cascadeStatus(in.Id, 0, tx); err != nil {
				l.Errorf("批量设置级联子孙关闭状态失败%v", err)
				return err
			}
		}
		// 1.2 开启级联子孙
		if in.HasStatus && in.Status == 1 && oldMenu.Status == 0 {
			if err := l.cascadeStatus(in.Id, 1, tx); err != nil {
				l.Errorf("批量设置级联子孙开启状态失败%v", err)
				return err
			}
		}

		// 2.1 隐藏级联子孙
		if in.HasVisible && in.Visible == 0 && oldMenu.Visible == 1 && oldMenu.MenuType != 2 {
			if err := l.cascadeVisible(in.Id, 0, tx); err != nil {
				l.Errorf("批量隐藏级联子孙失败%v", err)
				return err
			}
		}
		// 2.2 可见级联子孙
		if in.HasVisible && in.Visible == 1 && oldMenu.Visible == 0 && oldMenu.MenuType != 2 {
			if err := l.cascadeVisible(in.Id, 1, tx); err != nil {
				l.Errorf("批量可见级联子孙失败%v", err)
				return err
			}
		}

		if err := l.svcCtx.SysMenuModel.UpdateMenuTrans(l.ctx, tx, in.Id, updates); err != nil {
			l.Errorf("更新菜单[%d]失败：%v", in.Id, err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		return nil
	})

	if err != nil {
		l.Errorf("更新菜单事务执行失败: %v", err)
		return nil, err
	}

	return &sys.Empty{}, nil
}

// cascadeStatus 批量更新指定菜单所有子孙的状态。
func (l *UpdateMenuLogic) cascadeStatus(parentId, status int64, tx *gorm.DB) error {
	menus, err := l.svcCtx.SysMenuModel.ListAll(l.ctx)
	if err != nil {
		l.Errorf("查询全部菜单失败：%v", err)
		return err
	}

	menuIds := collectMenuChildren(menus, parentId)
	if len(menuIds) == 0 {
		return nil
	}
	if err := l.svcCtx.SysMenuModel.BatchUpdateMenuStatus(l.ctx, tx, menuIds, int(status)); err != nil {
		return err
	}
	return nil
}

// cascadeVisible 批量更新指定菜单所有子孙的可见性。
func (l *UpdateMenuLogic) cascadeVisible(parentId, visible int64, tx *gorm.DB) error {
	menus, err := l.svcCtx.SysMenuModel.ListAll(l.ctx)
	if err != nil {
		l.Errorf("查询全部菜单失败：%v", err)
		return err
	}

	menuIds := collectMenuChildren(menus, parentId)
	if len(menuIds) == 0 {
		return nil
	}
	if err := l.svcCtx.SysMenuModel.BatchUpdateMenuVisible(l.ctx, tx, menuIds, int(visible)); err != nil {
		return err
	}
	return nil
}

// collectMenuChildren 递归收集指定菜单的所有未软删除子孙ID。
func collectMenuChildren(menus []*systemmodel.SysMenu, parentId int64) []int64 {
	var ids []int64
	for _, m := range menus {
		if m.ParentId == parentId && !m.DeletedAt.Valid {
			ids = append(ids, m.Id)
			ids = append(ids, collectMenuChildren(menus, m.Id)...)
		}
	}
	return ids
}
