package common

import (
	"context"
	systemmodel "go-zero-rpc/sys-rpc/internal/model"
	casbinpkg "go-zero-rpc/sys-rpc/pkg/casbin"
	"go-zero-rpc/sys-rpc/sys"

	casbinv2 "github.com/casbin/casbin/v2"
	"github.com/zeromicro/go-zero/core/logx"
)

// RoleApiQuerier 提供根据角色ID查询绑定接口的能力。
type RoleApiQuerier interface {
	ListApiIdsByRoleId(ctx context.Context, roleId int64) ([]int64, error)
}

// ApiQuerier 提供根据接口ID批量查询接口详情的能力。
type ApiQuerier interface {
	ListByIds(ctx context.Context, ids []int64) ([]systemmodel.SysApi, error)
}

// RoleQuerier 提供根据角色ID批量查询角色信息的能力。
type RoleQuerier interface {
	FindByIds(ctx context.Context, roleIds []int64) ([]*systemmodel.SysRole, error)
}

// RebuildCasbinByRoleIds 根据角色ID列表，重建每个角色的Casbin策略。
//
// 适用场景：API的 path/method 变更、API被删除后，需要同步更新所有关联角色的Casbin规则。
// 该方法只记录错误日志，不返回错误（降级处理，不影响主流程）。
//
// 参数：
//   - ctx           : 上下文
//   - enforcer      : Casbin执行器
//   - roleQuerier   : 角色信息查询器（拿到 role.Code）
//   - roleApiQuerier: 角色-接口绑定查询器
//   - apiQuerier    : 接口详情查询器
//   - roleIds       : 需要重建Casbin策略的角色ID列表
func RebuildCasbinByRoleIds(
	ctx context.Context,
	enforcer *casbinv2.Enforcer,
	roleQuerier RoleQuerier,
	roleApiQuerier RoleApiQuerier,
	apiQuerier ApiQuerier,
	roleIds []int64,
) {
	if len(roleIds) == 0 {
		logx.Info("roleIds is empty")
		return
	}

	logger := logx.WithContext(ctx)

	// 1.批量查询角色信息（拿roleCode）
	roles, err := roleQuerier.FindByIds(ctx, roleIds)
	if err != nil {
		logger.Errorf("重建Casbin策略-查询角色信息失败：%v", err)
		return
	}

	// 2.逐个角色重建策略
	for _, role := range roles {
		// 2.1 查询该角色绑定的所有接口ID
		apiIds, err := roleApiQuerier.ListApiIdsByRoleId(ctx, role.Id)
		if err != nil {
			logger.Errorf("重建Casbin策略-查询角色[%s]绑定接口失败：%v", role.Code, err)
			continue
		}

		// 2.2 该角色没有绑定任何接口，清空策略
		if len(apiIds) == 0 {
			if err := casbinpkg.RemoveAllPoliciesForRole(enforcer, role.Code); err != nil {
				logger.Errorf("重建Casbin策略-清除角色[%s]策略失败：%v", role.Code, err)
			}
			continue
		}

		// 2.3 查询接口详情（拿path和method）
		apis, err := apiQuerier.ListByIds(ctx, apiIds)
		if err != nil {
			logger.Errorf("重建Casbin策略-查询接口列表失败：%v", err)
			continue
		}

		// 2.4 构建规则并全量覆盖
		rules := make([][]string, 0, len(apis))
		for _, api := range apis {
			rules = append(rules, []string{api.ApiPath, api.Method})
		}
		if err := casbinpkg.AddRolePolicies(enforcer, role.Code, rules); err != nil {
			logger.Errorf("重建Casbin策略-角色[%s]写入失败：%v", role.Code, err)
		}
	}
}

// MenuQuerier 菜单查询接口，用于解耦对 SysMenuModel 的依赖
type MenuQuerier interface {
	ListByIds(ctx context.Context, ids []int64) ([]*systemmodel.SysMenu, error)
}

// buildMenuTree 将扁平菜单列表递归构建为树形结构。
func BuildMenuTree(menus []*systemmodel.SysMenu, parentId int64) []*sys.MenuItem {
	result := make([]*sys.MenuItem, 0)
	for _, m := range menus {
		if m.ParentId != parentId {
			continue
		}
		menuNode := sys.MenuItem{
			Id:        m.Id,
			ParentId:  m.ParentId,
			MenuName:  m.Name,
			Path:      m.MenuPath,
			Component: m.Component,
			Icon:      m.Icon,
			MenuType:  m.MenuType,
			Perms:     m.Permission,
			Sort:      m.Sort,
			Status:    m.Status,
			Visible:   m.Visible,
		}
		children := BuildMenuTree(menus, m.Id)
		if len(children) > 0 {
			menuNode.Children = children
		}
		result = append(result, &menuNode)
	}
	return result
}

// 补全菜单祖先链辅助函数
// 第二个参数需要外面传一个能查菜单的对象

func CompleteMenuAncestors(ctx context.Context, querier MenuQuerier, menuIds []int64) ([]int64, error) {
	if len(menuIds) == 0 {
		return []int64{}, nil
	}

	// 用 map 去重，同时记录已经收集到的所有菜单ID
	allIds := make(map[int64]struct{}, len(menuIds))
	for _, id := range menuIds {
		allIds[id] = struct{}{}
	}

	// 第一轮：查询前端传入的菜单，收集需要补查的父ID
	toQuery := make([]int64, 0, len(menuIds))
	menus, err := querier.ListByIds(ctx, menuIds)
	if err != nil {
		return nil, err
	}
	for _, menu := range menus {
		if menu.ParentId == 0 {
			continue
		}
		if _, exists := allIds[menu.ParentId]; !exists {
			allIds[menu.ParentId] = struct{}{}
			toQuery = append(toQuery, menu.ParentId)
		}
	}

	// 循环向上查找，直到没有新的父节点需要补全
	for len(toQuery) > 0 {
		menus, err = querier.ListByIds(ctx, toQuery)
		if err != nil {
			return nil, err
		}

		toQuery = toQuery[:0] // 清空，准备收集下一轮
		for _, menu := range menus {
			if menu.ParentId == 0 {
				continue
			}
			if _, exists := allIds[menu.ParentId]; !exists {
				allIds[menu.ParentId] = struct{}{}
				toQuery = append(toQuery, menu.ParentId)
			}
		}
	}

	// map 转 slice
	result := make([]int64, 0, len(allIds))
	for id := range allIds {
		result = append(result, id)
	}
	return result, nil
}
