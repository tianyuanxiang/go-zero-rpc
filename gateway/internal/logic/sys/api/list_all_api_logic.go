// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package api

import (
	"context"

	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	sys "go-zero-rpc/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAllApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAllApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAllApiLogic {
	return &ListAllApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListAllApiLogic) ListAllApi() (resp *types.ListAllApiResp, err error) {
	rpcResp, err := l.svcCtx.SysRpc.ListAllApi(l.ctx, &sys.Empty{})
	if err != nil {
		l.Logger.Errorf("Call the ListAllApi failed, err=%v", err)
		return nil, err
	}

	groups := make([]types.ApiGroupOption, 0, len(rpcResp.Groups))
	for _, group := range rpcResp.Groups {
		if group == nil {
			continue
		}
		apis := make([]types.ApiOption, 0, len(group.Apis))
		for _, item := range group.Apis {
			if item == nil {
				continue
			}
			apis = append(apis, types.ApiOption{
				Id:      item.Id,
				ApiName: item.ApiName,
				ApiPath: item.ApiPath,
				Method:  item.Method,
				Remark:  item.Remark,
			})
		}
		groups = append(groups, types.ApiGroupOption{
			Group: group.Group,
			Apis:  apis,
		})
	}

	return &types.ListAllApiResp{
		Groups: groups,
	}, nil
}
