// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"
	"go-zero-rpc/sys-rpc/pb"

	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type GetRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleLogic {
	return &GetRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetRole 根据角色ID查询角色详情，包含关联的菜单ID和接口ID。
func (l *GetRoleLogic) GetRole(in *pb.GetRoleReq) (*pb.RoleItem, error) {
	role, err := l.svcCtx.SysRoleModel.FindOneByRoleId(l.ctx, in.RoleId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, xerr.NewCodeError(xerr.ErrRoleNotFound)
		}
		l.Errorf("查询角色[%d]失败：%v", in.RoleId, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	// 查询角色关联的菜单ID列表
	menuIds, err := l.svcCtx.SysRoleMenuModel.GetMenuIdsByRoleIds(l.ctx, []int64{in.RoleId})
	if err != nil {
		l.Errorf("查询角色[%d]关联菜单失败：%v", in.RoleId, err)
		menuIds = []int64{}
	}

	// 查询角色关联的接口ID列表
	apiIds, err := l.svcCtx.SysRoleApiModel.ListApiIdsByRoleId(l.ctx, in.RoleId)
	if err != nil {
		l.Errorf("查询角色[%d]关联接口失败：%v", in.RoleId, err)
		apiIds = []int64{}
	}

	return &pb.RoleItem{
		Id:        role.Id,
		RoleCode:  role.Code,
		RoleName:  role.Name,
		Status:    role.Status,
		Sort:      role.Sort,
		Remark:    role.Remark,
		CreatedAt: role.CreatedAt.Format("2006-01-02 15:04:05"),
		MenuIds:   menuIds,
		ApiIds:    apiIds,
	}, nil
}
