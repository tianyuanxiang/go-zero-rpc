// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package file

import (
	"context"

	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFileLogic {
	return &ListFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListFileLogic) ListFile(req *types.ListFileReq) (resp *types.ListFileResp, err error) {
	// todo: add your logic here and delete this line

	return
}
