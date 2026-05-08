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

type ListApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListApiLogic {
	return &ListApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListApiLogic) ListApi(req *types.ListApiReq) (resp *types.ListApiResp, err error) {
	rpcResp, err := l.svcCtx.SysRpc.ListApi(l.ctx, &sys.ListApiReq{
		Page:     int64(req.Page),
		PageSize: int64(req.PageSize),
		Keyword:  req.Keyword,
		Group:    req.Group,
	})
	if err != nil {
		l.Logger.Errorf("调用ListApi RPC失败, err=%v", err)
		return nil, err
	}

	list := make([]types.ApiItem, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		if item == nil {
			continue
		}
		list = append(list, types.ApiItem{
			Id:        item.Id,
			ApiName:   item.ApiName,
			ApiPath:   item.ApiPath,
			Method:    item.Method,
			Group:     item.Group,
			Remark:    item.Remark,
			CreatedAt: item.CreatedAt,
		})
	}

	return &types.ListApiResp{
		Total: rpcResp.Total,
		List:  list,
	}, nil
}
