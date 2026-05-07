// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dict

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-rpc/gateway/internal/svc"
)

type DeleteDictTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDictTypeLogic {
	return &DeleteDictTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteDictTypeLogic) DeleteDictType() error {
	// todo: add your logic here and delete this line

	return nil
}
