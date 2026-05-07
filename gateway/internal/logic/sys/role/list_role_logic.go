// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role

import (
	"context"

	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRoleLogic {
	return &ListRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRoleLogic) ListRole(req *types.ListRoleReq) (resp *types.ListRoleResp, err error) {
	// todo: add your logic here and delete this line

	return
}
