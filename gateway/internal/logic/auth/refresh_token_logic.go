// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"
	sys "go-zero-rpc/sys-rpc/pb"

	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshTokenLogic) RefreshToken(req *types.RefreshTokenReq) (resp *types.RefreshTokenResp, err error) {
	result, err := l.svcCtx.AuthRpc.RefreshToken(l.ctx, &sys.RefreshTokenReq{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		l.Logger.Errorf("Call the RefreshToken failed, err=%v", err)
		return nil, err
	}

	return &types.RefreshTokenResp{
		AccessToken: result.AccessToken,
		ExpiresIn:   result.ExpiresIn,
	}, err

}
