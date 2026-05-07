// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"
	"go-zero-rpc/sys-rpc/sys"

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
		l.Logger.Errorf("调用RefreshToken RPC失败", err)
		return nil, err
	}

	return &types.RefreshTokenResp{
		AccessToken: result.AccessToken,
		ExpiresIn:   result.ExpiresIn,
	}, err

}
