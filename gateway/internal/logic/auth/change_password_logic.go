// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"
	"go-zero-rpc/gateway/internal/middleware"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangePasswordLogic) ChangePassword(req *types.ChangePasswordReq) error {
	// 从JWT认证中间件注入的Context中获取当前登录用户ID
	userId := middleware.GetUserIdFromCtx(l.ctx)
	if userId == 0 {
		l.Errorf("userId is empty")
		return nil
	}

	_, err := l.svcCtx.AuthRpc.ChangePassword(l.ctx, &sys.ChangePasswordReq{
		UserId:      userId,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		return err
	}

	return err
}
