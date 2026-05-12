// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dict

import (
	"context"

	"go-zero-rpc/common/middleware"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteDictTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDictTypeLogic {
	return &DeleteDictTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteDictTypeLogic) DeleteDictType(dictTypeId int64) error {
	userId := middleware.GetUserIdFromCtx(l.ctx)

	_, err := l.svcCtx.SysRpc.DeleteDictType(l.ctx, &sys.DeleteDictTypeReq{
		DictTypeId: dictTypeId,
		OperatorId: userId,
	})
	if err != nil {
		l.Logger.Errorf("调用DeleteDictType RPC失败, operatorId=%d, dictTypeId=%d, err=%v", userId, dictTypeId, err)
		return err
	}

	return nil
}
