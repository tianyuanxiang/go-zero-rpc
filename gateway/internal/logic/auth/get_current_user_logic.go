// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"context"
	"go-zero-rpc/gateway/internal/middleware"
	"go-zero-rpc/sys-rpc/sys"

	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCurrentUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCurrentUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrentUserLogic {
	return &GetCurrentUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCurrentUserLogic) GetCurrentUser() (resp *types.UserInfoResp, err error) {
	// 从JWT认证中间件注入的Context中获取当前登录用户ID
	userId := middleware.GetUserIdFromCtx(l.ctx)
	if userId == 0 {
		l.Errorf("userId is empty")
		return nil, nil
	}

	rpcResp, err := l.svcCtx.AuthRpc.GetCurrentUser(l.ctx, &sys.GetCurrentUserReq{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	return &types.UserInfoResp{
		UserInfo: types.UserInfo{
			UserId:   rpcResp.UserInfo.UserId,
			Username: rpcResp.UserInfo.Username,
			Nickname: rpcResp.UserInfo.Nickname,
			Email:    rpcResp.UserInfo.Email,
			Phone:    rpcResp.UserInfo.Phone,
			Avatar:   rpcResp.UserInfo.Avatar,
			Roles:    rpcResp.UserInfo.Roles,
			Menus:    convertMenuItems(rpcResp.UserInfo.Menus),
		},
	}, err
}
