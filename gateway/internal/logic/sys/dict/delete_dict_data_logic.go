// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dict

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-rpc/gateway/internal/svc"
)

type DeleteDictDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteDictDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDictDataLogic {
	return &DeleteDictDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteDictDataLogic) DeleteDictData() error {
	// todo: add your logic here and delete this line

	return nil
}
