package common

import (
	"context"
	systemmodel "go-zero-rpc/sys-rpc/internal/model"
	"go-zero-rpc/sys-rpc/sys"
)

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
