// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dict

import (
	"context"

	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateDictTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDictTypeLogic {
	return &UpdateDictTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateDictTypeLogic) UpdateDictType(req *types.UpdateDictTypeReq) error {
	// todo: add your logic here and delete this line

	return nil
}
