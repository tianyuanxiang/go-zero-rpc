// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package api

import (
	"context"

	"go-zero-rpc/gateway/internal/middleware"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateApiLogic {
	return &UpdateApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateApiLogic) UpdateApi(req *types.UpdateApiReq) error {
	operatorId := middleware.GetUserIdFromCtx(l.ctx)

	rpcReq := &sys.UpdateApiReq{
		Id:         req.Id,
		OperatorId: operatorId,
	}

	if req.ApiName != nil {
		rpcReq.HasApiName = true
		rpcReq.ApiName = *req.ApiName
	}
	if req.ApiPath != nil {
		rpcReq.HasApiPath = true
		rpcReq.ApiPath = *req.ApiPath
	}
	if req.Method != nil {
		rpcReq.HasMethod = true
		rpcReq.Method = *req.Method
	}
	if req.Group != nil {
		rpcReq.HasGroup = true
		rpcReq.Group = *req.Group
	}
	if req.Remark != nil {
		rpcReq.HasRemark = true
		rpcReq.Remark = *req.Remark
	}

	_, err := l.svcCtx.SysRpc.UpdateApi(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("调用UpdateApi RPC失败, operatorId=%d, apiId=%d, err=%v", operatorId, req.Id, err)
		return err
	}

	return nil
}
