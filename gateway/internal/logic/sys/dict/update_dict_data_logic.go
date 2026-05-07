// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dict

import (
	"context"

	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateDictDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateDictDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDictDataLogic {
	return &UpdateDictDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDictDataLogic) UpdateDictData(req *types.UpdateDictDataReq) error {
	// todo: add your logic here and delete this line

	return nil
}
