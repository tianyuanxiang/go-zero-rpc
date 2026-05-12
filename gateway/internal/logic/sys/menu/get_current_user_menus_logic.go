// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package menu

import (
	"context"

	"go-zero-rpc/common/middleware"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCurrentUserMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCurrentUserMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrentUserMenusLogic {
	return &GetCurrentUserMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCurrentUserMenusLogic) GetCurrentUserMenus() (resp *types.MenuTreeResp, err error) {
	userId := middleware.GetUserIdFromCtx(l.ctx)
	if userId == 0 {
		l.Errorf("userId is empty")
		return nil, nil
	}

	rpcResp, err := l.svcCtx.SysRpc.GetCurrentUserMenus(l.ctx, &sys.GetCurrentUserMenusReq{
		UserId: userId,
	})
	if err != nil {
		l.Logger.Errorf("调用GetCurrentUserMenus RPC失败, userId=%d, err=%v", userId, err)
		return nil, err
	}

	return &types.MenuTreeResp{
		List: convertMenuItems(rpcResp.List),
	}, nil
}
