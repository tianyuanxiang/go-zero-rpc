// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package menu

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-rpc/gateway/internal/svc"
)

type DeleteMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuLogic {
	return &DeleteMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMenuLogic) DeleteMenu() error {
	// todo: add your logic here and delete this line

	return nil
}
