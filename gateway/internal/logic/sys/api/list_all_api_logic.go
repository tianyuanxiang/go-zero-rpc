// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package api

import (
	"context"

	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	"go-zero-rpc/sys-rpc/sys"

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
		l.Logger.Errorf("调用ListAllApi RPC失败, err=%v", err)
		return nil, err
	}

	listAll := make([]types.ApiOption, 0, len(rpcResp.ListAll))
	for _, item := range rpcResp.ListAll {
		if item == nil {
			continue
		}
		listAll = append(listAll, types.ApiOption{
			Id:      item.Id,
			ApiName: item.ApiName,
			ApiPath: item.ApiPath,
			Remark:  item.Remark,
		})
	}

	return &types.ListAllApiResp{
		ListAll: listAll,
	}, nil
}
