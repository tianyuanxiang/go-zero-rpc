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

type CreateDictTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDictTypeLogic {
	return &CreateDictTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateDictTypeLogic) CreateDictType(req *types.CreateDictTypeReq) error {
	operatorId := middleware.GetUserIdFromCtx(l.ctx)

	_, err := l.svcCtx.SysRpc.CreateDictType(l.ctx, &sys.CreateDictTypeReq{
		DictName:   req.DictName,
		DictCode:   req.DictCode,
		Remark:     req.Remark,
		Status:     int64(req.Status),
		OperatorId: operatorId,
	})
	if err != nil {
		l.Logger.Errorf("调用CreateDictType RPC失败, operatorId=%d, err=%v", operatorId, err)
		return err
	}

	return nil
}
