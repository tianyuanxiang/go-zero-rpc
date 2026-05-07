package user

import (
	"context"

	"go-zero-rpc/gateway/internal/middleware"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.DeleteUserReq) error {
	userId := middleware.GetUserIdFromCtx(l.ctx)
	if userId == 0 {
		return nil
	}

	_, err := l.svcCtx.SysRpc.DeleteUser(l.ctx, &sys.DeleteUserReq{
		UserId:     req.Id,
		OperatorId: userId,
	})
	if err != nil {
		l.Logger.Errorf("调用DeleteUser RPC失败, operatorId=%d, targetId=%d, err=%v", userId, req.Id, err)
		return err
	}

	return nil
}
