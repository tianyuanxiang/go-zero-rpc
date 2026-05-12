// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dict

import (
	"context"

	"go-zero-rpc/common/middleware"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateDictTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDictTypeLogic {
	return &UpdateDictTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDictTypeLogic) UpdateDictType(req *types.UpdateDictTypeReq) error {
	operatorId := middleware.GetUserIdFromCtx(l.ctx)

	rpcReq := &sys.UpdateDictTypeReq{
		Id:         int64(req.Id),
		OperatorId: operatorId,
	}

	if req.DictName != nil {
		rpcReq.HasDictName = true
		rpcReq.DictName = *req.DictName
	}
	if req.DictCode != nil {
		rpcReq.HasDictCode = true
		rpcReq.DictCode = *req.DictCode
	}
	if req.Remark != nil {
		rpcReq.HasRemark = true
		rpcReq.Remark = *req.Remark
	}
	if req.Status != nil {
		rpcReq.HasStatus = true
		rpcReq.Status = int64(*req.Status)
	}

	_, err := l.svcCtx.SysRpc.UpdateDictType(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("调用UpdateDictType RPC失败, operatorId=%d, dictTypeId=%d, err=%v", operatorId, req.Id, err)
		return err
	}

	return nil
}
