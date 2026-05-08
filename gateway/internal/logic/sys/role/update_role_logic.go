// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role

import (
	"context"

	"go-zero-rpc/gateway/internal/middleware"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRoleLogic) UpdateRole(req *types.UpdateRoleReq) error {
	operatorId := middleware.GetUserIdFromCtx(l.ctx)

	rpcReq := &sys.UpdateRoleReq{
		Id:         req.Id,
		OperatorId: operatorId,
	}

	if req.RoleCode != nil {
		rpcReq.HasRoleCode = true
		rpcReq.RoleCode = *req.RoleCode
	}
	if req.RoleName != nil {
		rpcReq.HasRoleName = true
		rpcReq.RoleName = *req.RoleName
	}
	if req.Sort != nil {
		rpcReq.HasSort = true
		rpcReq.Sort = *req.Sort
	}
	if req.Remark != nil {
		rpcReq.HasRemark = true
		rpcReq.Remark = *req.Remark
	}

	_, err := l.svcCtx.SysRpc.UpdateRole(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("调用UpdateRole RPC失败, operatorId=%d, roleId=%d, err=%v", operatorId, req.Id, err)
		return err
	}

	return nil
}
