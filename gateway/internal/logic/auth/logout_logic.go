// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"
	"go-zero-rpc/common/middleware"
	sys "go-zero-rpc/sys-rpc/pb"
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
	// 浠嶫WT璁よ瘉涓棿浠舵敞鍏ョ殑Context涓幏鍙栧綋鍓嶇櫥褰曠敤鎴稩D
	userId := middleware.GetUserIdFromCtx(l.ctx)
	if userId == 0 {
		l.Errorf("userId is empty")
		return nil
	}

	// 浠?Authorization header 鎻愬彇 token
	authHeader := r.Header.Get("Authorization")
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	_, err := l.svcCtx.AuthRpc.Logout(l.ctx, &sys.LogoutReq{
		AccessToken:               tokenStr,
		UserId:                    userId,
		RedisTokenBlacklistPrefix: middleware.RedisTokenBlacklistPrefix,
	})
	if err != nil {
		l.Logger.Errorf("Call the Logout failed, userId=%d, err=%v", userId, err)
		return err
	}

	return nil
}
