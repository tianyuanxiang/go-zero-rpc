// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role

import (
	"context"

	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	"go-zero-rpc/sys-rpc/sys"

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
	rpcResp, err := l.svcCtx.SysRpc.ListRole(l.ctx, &sys.ListRoleReq{
		Page:     int64(req.Page),
		PageSize: int64(req.PageSize),
		Keyword:  req.Keyword,
	})
	if err != nil {
		l.Logger.Errorf("调用ListRole RPC失败, err=%v", err)
		return nil, err
	}

	list := make([]types.RoleItem, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		if item == nil {
			continue
		}
		list = append(list, types.RoleItem{
			Id:        item.Id,
			RoleName:  item.RoleName,
			RoleCode:  item.RoleCode,
			Status:    int(item.Status),
			Sort:      int(item.Sort),
			Remark:    item.Remark,
			CreatedAt: item.CreatedAt,
		})
	}

	return &types.ListRoleResp{
		Total: rpcResp.Total,
		List:  list,
	}, nil
}
