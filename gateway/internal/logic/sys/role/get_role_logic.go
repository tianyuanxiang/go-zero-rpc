// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role

import (
	"context"

	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleLogic {
	return &GetRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRoleLogic) GetRole(roleId int64) (resp *types.RoleItem, err error) {
	role, err := l.svcCtx.SysRpc.GetRole(l.ctx, &sys.GetRoleReq{
		RoleId: roleId,
	})
	if err != nil {
		l.Logger.Errorf("调用GetRole RPC失败, roleId=%d, err=%v", roleId, err)
		return nil, err
	}

	return &types.RoleItem{
		Id:        role.Id,
		RoleName:  role.RoleName,
		RoleCode:  role.RoleCode,
		Status:    int(role.Status),
		Sort:      int(role.Sort),
		Remark:    role.Remark,
		CreatedAt: role.CreatedAt,
	}, nil
}
