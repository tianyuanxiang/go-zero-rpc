// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dict

import (
	"context"

	"go-zero-rpc/common/middleware"
	"go-zero-rpc/gateway/internal/svc"
	sys "go-zero-rpc/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteDictDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteDictDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDictDataLogic {
	return &DeleteDictDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteDictDataLogic) DeleteDictData(dictDataId int64) error {
	userId := middleware.GetUserIdFromCtx(l.ctx)

	_, err := l.svcCtx.SysRpc.DeleteDictData(l.ctx, &sys.DeleteDictDataReq{
		DictDataId: dictDataId,
		OperatorId: userId,
	})
	if err != nil {
		l.Logger.Errorf("Call the DeleteDictData failed, operatorId=%d, dictDataId=%d, err=%v", userId, dictDataId, err)
		return err
	}

	return nil
}
