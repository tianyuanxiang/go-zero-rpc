// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"
	"go-zero-rpc/sys-rpc/pb"

	"go-zero-rpc/common/constants"
	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/common"
	systemmodel "go-zero-rpc/sys-rpc/internal/model"
	"go-zero-rpc/sys-rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCurrentUserMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCurrentUserMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrentUserMenusLogic {
	return &GetCurrentUserMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetCurrentUserMenus 根据用户ID返回该用户可见的菜单树。
//
// 业务流程：
//  1. 查询用户所有角色ID
//  2. 若是 admin，则返回全部菜单；否则按角色绑定的菜单查询
//  3. 非 admin 的菜单需补全祖先链，确保父菜单不会缺失
//  4. 过滤隐藏菜单
//  5. 构建菜单树
func (l *GetCurrentUserMenusLogic) GetCurrentUserMenus(in *pb.GetCurrentUserMenusReq) (*pb.MenuTreeResp, error) {
	userId := in.UserId
	if userId == 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	roleIds, err := l.svcCtx.SysUserRoleModel.GetRoleIdsByUserId(l.ctx, userId)
	if err != nil {
		l.Errorf("查询用户[%d]角色失败：%v", userId, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	if len(roleIds) == 0 {
		l.Errorf("用户[%d]角色ID为空", userId)
		return &pb.MenuTreeResp{List: []*pb.MenuItem{}}, nil
	}

	var menus []*systemmodel.SysMenu
	if l.isAdmin(roleIds) {
		menus, err = l.svcCtx.SysMenuModel.ListAll(l.ctx)
		if err != nil {
			l.Errorf("查询超管所有菜单失败：%v", err)
			return nil, xerr.NewCodeError(xerr.ErrInternal)
		}
	} else {
		menuIds, mErr := l.svcCtx.SysRoleMenuModel.GetMenuIdsByRoleIds(l.ctx, roleIds)
		if mErr != nil {
			l.Errorf("查询角色菜单权限失败：%v", mErr)
			return nil, xerr.NewCodeError(xerr.ErrInternal)
		}
		menuIds, err = common.CompleteMenuAncestors(l.ctx, l.svcCtx.SysMenuModel, menuIds)
		if err != nil {
			l.Errorf("补齐菜单祖先链失败：%v", err)
			return nil, xerr.NewCodeError(xerr.ErrInternal)
		}
		menus, err = l.svcCtx.SysMenuModel.ListByIds(l.ctx, menuIds)
		if err != nil {
			l.Errorf("查询菜单失败：%v", err)
			return nil, xerr.NewCodeError(xerr.ErrInternal)
		}
	}

	// 过滤隐藏菜单
	visible := make([]*systemmodel.SysMenu, 0, len(menus))
	for _, m := range menus {
		if m.Visible == 1 {
			visible = append(visible, m)
		}
	}

	menuTree := common.BuildMenuTree(visible, 0)
	return &pb.MenuTreeResp{List: menuTree}, nil
}

// isAdmin 判断给定角色ID列表中是否含有超级管理员角色。
func (l *GetCurrentUserMenusLogic) isAdmin(roleIds []int64) bool {
	for _, roleId := range roleIds {
		role, err := l.svcCtx.SysRoleModel.FindOneByRoleId(l.ctx, roleId)
		if err != nil || role == nil || role.Status != 1 {
			continue
		}
		if role.Code == constants.RoleCodeAdmin {
			return true
		}
	}
	return false
}
