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

type CreateDictDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateDictDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDictDataLogic {
	return &CreateDictDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateDictDataLogic) CreateDictData(req *types.CreateDictDataReq) error {
	operatorId := middleware.GetUserIdFromCtx(l.ctx)

	_, err := l.svcCtx.SysRpc.CreateDictData(l.ctx, &sys.CreateDictDataReq{
		DictTypeId: int64(req.DictTypeId),
		DictLabel:  req.DictLabel,
		DictValue:  req.DictValue,
		Sort:       int64(req.Sort),
		Remark:     req.Remark,
		Status:     int64(req.Status),
		OperatorId: operatorId,
	})
	if err != nil {
		l.Logger.Errorf("调用CreateDictData RPC失败, operatorId=%d, err=%v", operatorId, err)
		return err
	}

	return nil
}
