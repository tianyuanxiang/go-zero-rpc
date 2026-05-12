// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role

import (
	"context"

	"go-zero-rpc/common/middleware"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	"go-zero-rpc/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRolePermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRolePermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRolePermissionsLogic {
	return &UpdateRolePermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRolePermissionsLogic) UpdateRolePermissions(req *types.UpdateRolePermissionsReq) error {
	operatorId := middleware.GetUserIdFromCtx(l.ctx)

	_, err := l.svcCtx.SysRpc.UpdateRolePermissions(l.ctx, &sys.UpdateRolePermissionsReq{
		Id:         req.Id,
		MenuIds:    req.MenuIds,
		ApiIds:     req.ApiIds,
		OperatorId: operatorId,
	})
	if err != nil {
		l.Logger.Errorf("调用UpdateRolePermissions RPC失败, operatorId=%d, roleId=%d, err=%v", operatorId, req.Id, err)
		return err
	}

	return nil
}
