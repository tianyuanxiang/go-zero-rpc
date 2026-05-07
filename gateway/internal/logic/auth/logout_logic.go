// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"
	"go-zero-rpc/gateway/internal/middleware"
	"go-zero-rpc/sys-rpc/sys"
	"net/http"
	"strings"

	"go-zero-rpc/gateway/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(r *http.Request) error {
	// 从JWT认证中间件注入的Context中获取当前登录用户ID
	userId := middleware.GetUserIdFromCtx(l.ctx)
	if userId == 0 {
		l.Errorf("userId is empty")
		return nil
	}

	// 从 Authorization header 提取 token
	authHeader := r.Header.Get("Authorization")
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	_, err := l.svcCtx.AuthRpc.Logout(l.ctx, &sys.LogoutReq{
		AccessToken:               tokenStr,
		UserId:                    userId,
		RedisTokenBlacklistPrefix: middleware.RedisTokenBlacklistPrefix,
	})
	if err != nil {
		l.Logger.Errorf("调用Logout RPC失败, userId=%d, err=%v", userId, err)
		return err
	}

	return nil
}
