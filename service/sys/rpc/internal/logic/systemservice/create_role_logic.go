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

type CreateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateRole 创建新角色，并同步关联菜单、接口和Casbin策略。
//
// 业务流程：
//  1. 校验角色编码唯一性
//  2. 校验关联API与菜单存在性
//  3. 事务内插入角色、关联菜单、关联接口（菜单需补全祖先链）
//  4. 事务成功后同步Casbin策略
//
// 参数：
//   - in : 创建角色请求体
//
// 返回：
//   - *sys.CreateRoleResp : 携带新建角色ID
//   - error               : 业务错误
func (l *CreateRoleLogic) CreateRole(in *sys.CreateRoleReq) (*sys.CreateRoleResp, error) {
	// 1. 检查角色编码唯一性
	existRole, err := l.svcCtx.SysRoleModel.FindOneByCode(l.ctx, in.RoleCode)
	if err != nil && err != sqlx.ErrNotFound {
		l.Errorf("查询角色编码[%s]是否存在失败：%v", in.RoleCode, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}
	if existRole != nil {
		return nil, xerr.NewCodeError(xerr.ErrRoleCodeDuplicate)
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

	// 用闭包外变量捕获事务内生成的角色ID，事务成功后供 resp 返回
	var newRoleId int64
	menuIds := append([]int64{}, in.MenuIds...)

	// 4. 开启事务
	err = l.svcCtx.Orm.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		// 4.1 插入角色
		roleId, err := l.svcCtx.SysRoleModel.InsertRoleTrans(l.ctx, tx, &systemmodel.SysRole{
			Name:   in.RoleName,
			Code:   in.RoleCode,
			Status: 1,
			Remark: in.Remark,
			Sort:   in.Sort,
		})
		if err != nil {
			l.Errorf("插入角色记录失败：%v", err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		newRoleId = roleId

		// 4.2 插入关联菜单ID（先补全祖先链）
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
					RoleId: roleId,
					MenuId: menuId,
				})
			}
			if _, err := l.svcCtx.SysRoleMenuModel.InsertRoleMenuTrans(l.ctx, tx, roleMenus); err != nil {
				l.Errorf("插入关联菜单ID失败：%v", err)
				return xerr.NewCodeError(xerr.ErrInternal)
			}
		}

		// 4.3 插入关联接口ID
		if len(in.ApiIds) > 0 {
			roleApis := make([]systemmodel.SysRoleApi, 0, len(in.ApiIds))
			for _, apiId := range in.ApiIds {
				roleApis = append(roleApis, systemmodel.SysRoleApi{
					RoleId: roleId,
					ApiId:  apiId,
				})
			}
			if _, err := l.svcCtx.SysRoleApiModel.InsertRoleApiTrans(l.ctx, tx, roleApis); err != nil {
				l.Errorf("插入关联接口ID失败：%v", err)
				return xerr.NewCodeError(xerr.ErrInternal)
			}
		}

		return nil
	})
	if err != nil {
		l.Errorf("创建角色事务执行失败: %v", err)
		return nil, err
	}

	// 5. 同步Casbin策略
	if len(apiModels) > 0 {
		rules := make([][]string, 0, len(apiModels))
		for _, api := range apiModels {
			rules = append(rules, []string{api.ApiPath, api.Method})
		}
		if err := casbinpkg.AddRolePolicies(l.svcCtx.Enforcer, in.RoleCode, rules); err != nil {
			l.Errorf("同步Casbin策略失败（角色编码: %s）：%v", in.RoleCode, err)
		}
	}

	return &sys.CreateRoleResp{RoleId: newRoleId}, nil
}
