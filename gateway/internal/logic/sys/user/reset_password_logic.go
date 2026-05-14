package user

import (
	"context"

	"go-zero-rpc/common/middleware"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	sys "go-zero-rpc/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetPasswordLogic) ResetPassword(req *types.ResetPasswordReq) error {
	userId := middleware.GetUserIdFromCtx(l.ctx)
	if userId == 0 {
		return nil
	}

	_, err := l.svcCtx.SysRpc.ResetPassword(l.ctx, &sys.ResetPasswordReq{
		UserId:      req.Id,
		NewPassword: req.NewPassword,
		OperatorId:  userId,
	})
	if err != nil {
		l.Logger.Errorf("Call the ResetPassword failed, operatorId=%d, targetId=%d, err=%v", userId, req.Id, err)
		return err
	}

	return nil
}
