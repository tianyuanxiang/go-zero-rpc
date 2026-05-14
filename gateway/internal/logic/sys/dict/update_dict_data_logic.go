// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dict

import (
	"context"

	"go-zero-rpc/common/middleware"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	sys "go-zero-rpc/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateDictDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateDictDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDictDataLogic {
	return &UpdateDictDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDictDataLogic) UpdateDictData(req *types.UpdateDictDataReq) error {
	operatorId := middleware.GetUserIdFromCtx(l.ctx)

	rpcReq := &sys.UpdateDictDataReq{
		Id:         int64(req.Id),
		DictTypeId: int64(req.DictTypeId),
		DictLabel:  req.DictLabel,
		DictValue:  req.DictValue,
		OperatorId: operatorId,
	}

	if req.Sort != nil {
		rpcReq.HasSort = true
		rpcReq.Sort = int64(*req.Sort)
	}
	if req.Remark != nil {
		rpcReq.HasRemark = true
		rpcReq.Remark = *req.Remark
	}
	if req.Status != nil {
		rpcReq.HasStatus = true
		rpcReq.Status = int64(*req.Status)
	}

	_, err := l.svcCtx.SysRpc.UpdateDictData(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Call the UpdateDictData failed, operatorId=%d, dictDataId=%d, err=%v", operatorId, req.Id, err)
		return err
	}

	return nil
}
