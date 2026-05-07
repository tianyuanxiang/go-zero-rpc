package authservicelogic

import (
	"context"
	"fmt"
	"go-zero-rpc/common/jwtx"
	"go-zero-rpc/common/xerr"
	"time"

	"go-zero-rpc/sys-rpc/internal/svc"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LogoutLogic) Logout(in *sys.LogoutReq) (*sys.CommonResp, error) {
	tokenStr := in.AccessToken
	if tokenStr == "" {
		l.Errorf("logout err: token empty")
		return nil, xerr.NewCodeError(xerr.ErrTokenInvalid)
	}

	// 解析Token获取过期时间（即使Token即将过期也要加入黑名单）
	expireAt := jwtx.GetTokenExpireAt(tokenStr, l.svcCtx.Config.JwtAuth.AccessSecret)
	if expireAt == 0 {
		// Token已过期或无效，无需加入黑名单
		return nil, xerr.NewCodeError(xerr.ErrTokenExpired)
	}

	// 计算剩余有效时间（秒）
	remainDuration := time.Until(time.Unix(expireAt, 0))
	if remainDuration <= 0 {
		// 已过期，无需处理
		return nil, xerr.NewCodeError(xerr.ErrTokenExpired)
	}

	// 将Token加入Redis黑名单（键：token:blacklist:<token>，TTL=剩余有效期）
	blacklistKey := fmt.Sprintf("%s%s", in.RedisTokenBlacklistPrefix, tokenStr)
	if err := l.svcCtx.RDB.Set(l.ctx, blacklistKey, "1", remainDuration).Err(); err != nil {
		l.Logger.Errorf("Token加入Redis黑名单失败：%v", err)
		// 黑名单写入失败不影响登出响应，日志记录即可
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}
	l.Infof("用户[%d]登出成功，Token已加入黑名单", in.UserId)

	return &sys.CommonResp{
		Code: 0,
		Msg:  "已登出",
	}, nil
}
