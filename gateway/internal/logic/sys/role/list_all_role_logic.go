// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role

import (
	"context"

	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	sys "go-zero-rpc/sys-rpc/pb"

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
	rpcResp, err := l.svcCtx.SysRpc.ListAllRole(l.ctx, &sys.Empty{})
	if err != nil {
		l.Logger.Errorf("Call the ListAllRole failed, err=%v", err)
		return nil, err
	}

	listAll := make([]types.RoleOption, 0, len(rpcResp.ListAll))
	for _, item := range rpcResp.ListAll {
		if item == nil {
			continue
		}
		listAll = append(listAll, types.RoleOption{
			Id:       item.Id,
			RoleName: item.RoleName,
			RoleCode: item.RoleCode,
		})
	}

	return &types.ListAllResp{
		ListAll: listAll,
	}, nil
}
