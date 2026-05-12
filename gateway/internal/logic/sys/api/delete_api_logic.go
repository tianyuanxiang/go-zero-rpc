// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package api

import (
	"context"

	"go-zero-rpc/common/middleware"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteApiLogic {
	return &DeleteApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteApiLogic) DeleteApi(apiId int64) error {
	userId := middleware.GetUserIdFromCtx(l.ctx)

	_, err := l.svcCtx.SysRpc.DeleteApi(l.ctx, &sys.DeleteApiReq{
		ApiId:      apiId,
		OperatorId: userId,
	})
	if err != nil {
		l.Logger.Errorf("调用DeleteApi RPC失败, operatorId=%d, apiId=%d, err=%v", userId, apiId, err)
		return err
	}

	return nil
}
