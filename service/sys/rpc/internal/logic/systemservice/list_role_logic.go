// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"

	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/svc"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRoleLogic {
	return &ListRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListRole 分页查询角色列表，支持关键词模糊检索。
func (l *ListRoleLogic) ListRole(in *sys.ListRoleReq) (*sys.ListRoleResp, error) {
	roles, count, err := l.svcCtx.SysRoleModel.List(l.ctx, int(in.Page), int(in.PageSize), in.Keyword)
	if err != nil {
		l.Errorf("查询角色列表失败 %v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	// 批量查询角色关联的菜单ID和接口ID
	roleIds := make([]int64, 0, len(roles))
	for _, role := range roles {
		roleIds = append(roleIds, role.Id)
	}

	// 构建角色菜单关联映射
	roleMenuMap := make(map[int64][]int64)
	roleApiMap := make(map[int64][]int64)
	if len(roleIds) > 0 {
		// 逐角色查询菜单ID
		for _, roleId := range roleIds {
			menuIds, mErr := l.svcCtx.SysRoleMenuModel.GetMenuIdsByRoleIds(l.ctx, []int64{roleId})
			if mErr == nil {
				roleMenuMap[roleId] = menuIds
			}
		}
		// 逐角色查询接口ID
		for _, roleId := range roleIds {
			apiIds, aErr := l.svcCtx.SysRoleApiModel.ListApiIdsByRoleId(l.ctx, roleId)
			if aErr == nil {
				roleApiMap[roleId] = apiIds
			}
		}
	}

	list := make([]*sys.RoleItem, 0, len(roles))
	for _, role := range roles {
		item := &sys.RoleItem{
			Id:        role.Id,
			RoleName:  role.Name,
			RoleCode:  role.Code,
			Status:    role.Status,
			Sort:      role.Sort,
			Remark:    role.Remark,
			CreatedAt: role.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if menuIds, ok := roleMenuMap[role.Id]; ok {
			item.MenuIds = menuIds
		} else {
			item.MenuIds = []int64{}
		}
		if apiIds, ok := roleApiMap[role.Id]; ok {
			item.ApiIds = apiIds
		} else {
			item.ApiIds = []int64{}
		}
		list = append(list, item)
	}

	return &sys.ListRoleResp{
		Total: count,
		List:  list,
	}, nil
}
