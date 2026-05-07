// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role

import (
	"context"

	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAllRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAllRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAllRoleLogic {
	return &ListAllRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAllRoleLogic) ListAllRole() (resp *types.ListAllResp, err error) {
	// todo: add your logic here and delete this line

	return
}
