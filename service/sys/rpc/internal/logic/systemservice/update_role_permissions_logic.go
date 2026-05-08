// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"

	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/common"
	systemmodel "go-zero-rpc/sys-rpc/internal/model"
	"go-zero-rpc/sys-rpc/internal/svc"
	casbinpkg "go-zero-rpc/sys-rpc/pkg/casbin"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

type UpdateRolePermissionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRolePermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRolePermissionsLogic {
	return &UpdateRolePermissionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateRolePermissions 更新角色的菜单与API权限关联，并同步Casbin策略。
//
// 业务流程：
//  1. 检查角色存在性
//  2. 校验关联API与菜单存在性
//  3. 事务内先删后插：更新菜单关联（菜单需补全祖先链）、更新接口关联
//  4. 事务成功后全量覆盖Casbin策略
func (l *UpdateRolePermissionsLogic) UpdateRolePermissions(in *sys.UpdateRolePermissionsReq) (*sys.Empty, error) {
	// 1. 检查角色是否存在
	existRole, err := l.svcCtx.SysRoleModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, xerr.NewCodeError(xerr.ErrRoleNotFound)
		}
		l.Errorf("查询角色[%d]失败：%v", in.Id, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}
	if existRole.DeletedAt.Valid {
		l.Errorf("角色[%s]已删除", existRole.Code)
		return nil, xerr.NewCodeError(xerr.ErrParamInvalid)
	}

	// 2. 获取接口信息并确认接口数量完整
	var apiModels []systemmodel.SysApi
	if len(in.ApiIds) > 0 {
		apiModels, err = l.svcCtx.SysApiModel.ListByIds(l.ctx, in.ApiIds)
		if err != nil {
			l.Errorf("查询接口信息失败：%v", err)
			return nil, xerr.NewCodeError(xerr.ErrInternal)
		}
		if len(apiModels) != len(in.ApiIds) {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "部分api接口信息不存在")
		}
	}

	// 3. 校验菜单存在性（含未软删除）
	if len(in.MenuIds) > 0 {
		menus, err := l.svcCtx.SysMenuModel.ListByIds(l.ctx, in.MenuIds)
		if err != nil {
			l.Errorf("查询菜单信息失败：%v", err)
			return nil, xerr.NewCodeError(xerr.ErrInternal)
		}
		if len(menus) != len(in.MenuIds) {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "部分菜单信息不存在或已删除")
		}
	}

	menuIds := append([]int64{}, in.MenuIds...)

	// 4. 开启事务
	err = l.svcCtx.Orm.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		// 4.1 先删后插：更新关联菜单
		if err := l.svcCtx.SysRoleMenuModel.DeleteRoleMenuByRoleIdTrans(l.ctx, tx, in.Id); err != nil {
			l.Errorf("删除角色[%d]旧菜单关联失败：%v", in.Id, err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		if len(menuIds) > 0 {
			completedIds, cErr := common.CompleteMenuAncestors(l.ctx, l.svcCtx.SysMenuModel, menuIds)
			if cErr != nil {
				l.Errorf("补全菜单祖先链失败：%v", cErr)
				return xerr.NewCodeError(xerr.ErrInternal)
			}
			menuIds = completedIds
			roleMenus := make([]systemmodel.SysRoleMenu, 0, len(menuIds))
			for _, menuId := range menuIds {
				roleMenus = append(roleMenus, systemmodel.SysRoleMenu{
					RoleId: in.Id,
					MenuId: menuId,
				})
			}
			if _, err := l.svcCtx.SysRoleMenuModel.InsertRoleMenuTrans(l.ctx, tx, roleMenus); err != nil {
				l.Errorf("插入角色[%d]新菜单关联失败：%v", in.Id, err)
				return xerr.NewCodeError(xerr.ErrInternal)
			}
		}

		// 4.2 先删后插：更新关联接口
		if err := l.svcCtx.SysRoleApiModel.DeleteRoleApiByRoleIdTrans(l.ctx, tx, in.Id); err != nil {
			l.Errorf("删除角色[%d]旧接口关联失败：%v", in.Id, err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		if len(in.ApiIds) > 0 {
			roleApis := make([]systemmodel.SysRoleApi, 0, len(in.ApiIds))
			for _, apiId := range in.ApiIds {
				roleApis = append(roleApis, systemmodel.SysRoleApi{
					RoleId: in.Id,
					ApiId:  apiId,
				})
			}
			if _, err := l.svcCtx.SysRoleApiModel.InsertRoleApiTrans(l.ctx, tx, roleApis); err != nil {
				l.Errorf("插入角色[%d]新接口关联失败：%v", in.Id, err)
				return xerr.NewCodeError(xerr.ErrInternal)
			}
		}
		return nil
	})

	if err != nil {
		l.Errorf("更新角色权限事务执行失败: %v", err)
		return nil, err
	}

	// 5. 同步Casbin策略（事务成功后执行，全量覆盖）
	if len(apiModels) > 0 {
		rules := make([][]string, 0, len(apiModels))
		for _, api := range apiModels {
			rules = append(rules, []string{api.ApiPath, api.Method})
		}
		if err := casbinpkg.AddRolePolicies(l.svcCtx.Enforcer, existRole.Code, rules); err != nil {
			l.Errorf("同步Casbin策略失败（角色编码: %s）：%v", existRole.Code, err)
		}
	} else {
		// API列表为空，清除该角色所有策略
		if err := casbinpkg.RemoveAllPoliciesForRole(l.svcCtx.Enforcer, existRole.Code); err != nil {
			l.Errorf("清除Casbin策略失败（角色编码: %s）：%v", existRole.Code, err)
		}
	}

	return &sys.Empty{}, nil
}
