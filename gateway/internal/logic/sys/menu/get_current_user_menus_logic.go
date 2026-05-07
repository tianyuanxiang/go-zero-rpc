// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package menu

import (
	"context"

	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
