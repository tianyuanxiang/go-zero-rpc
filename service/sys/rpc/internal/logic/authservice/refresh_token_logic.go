package authservicelogic

import (
	"context"
	"go-zero-rpc/common/jwtx"
	"go-zero-rpc/common/xerr"

	"go-zero-rpc/sys-rpc/internal/svc"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefreshTokenLogic) RefreshToken(in *sys.RefreshTokenReq) (*sys.RefreshTokenResp, error) {
	// 1. 解析刷新令牌
	claims, err := jwtx.ParseToken(in.RefreshToken, l.svcCtx.Config.JwtAuth.AccessSecret)
	if err != nil {
		l.Logger.Infof("刷新令牌验证失败：%v", err)
		return nil, xerr.NewCodeError(xerr.ErrTokenExpired)
	}

	// 2. 确认令牌类型必须是 RefreshToken
	if claims.TokenType != jwtx.TokenTypeRefresh {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrTokenInvalid, "请使用刷新令牌换取新的访问令牌")
	}

	// 3. 验证用户状态（防止用户被禁用后仍能刷新token）
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, claims.UserId)
	if err != nil {
		l.Logger.Errorf("查询用户[%d]失败：%v", claims.UserId, err)
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	if user.Status != 1 {
		return nil, xerr.NewCodeError(xerr.ErrAccountDisabled)
	}

	// 4. 生成新的访问令牌
	newAccessToken, err := jwtx.GenerateToken(user.Id, user.Username, l.svcCtx.Config.JwtAuth.AccessSecret, l.svcCtx.Config.JwtAuth.AccessExpire)
	if err != nil {
		l.Errorf("生成新AccessToken失败：%v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	return &sys.RefreshTokenResp{
		AccessToken: newAccessToken,
		ExpiresIn:   l.svcCtx.Config.JwtAuth.AccessExpire,
	}, err

}
