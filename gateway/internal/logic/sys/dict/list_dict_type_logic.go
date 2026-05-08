// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dict

import (
	"context"

	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDictTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDictTypeLogic {
	return &ListDictTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListDictTypeLogic) ListDictType(req *types.ListDictTypeReq) (resp *types.ListDictTypeResp, err error) {
	rpcResp, err := l.svcCtx.SysRpc.ListDictType(l.ctx, &sys.ListDictTypeReq{
		Page:     int64(req.Page),
		PageSize: int64(req.PageSize),
		Keyword:  req.Keyword,
	})
	if err != nil {
		l.Logger.Errorf("调用ListDictType RPC失败, err=%v", err)
		return nil, err
	}

	list := make([]types.DictTypeItem, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		if item == nil {
			continue
		}
		list = append(list, types.DictTypeItem{
			Id:        int(item.Id),
			DictName:  item.DictName,
			DictCode:  item.DictCode,
			Remark:    item.Remark,
			Status:    int(item.Status),
			CreatedAt: item.CreatedAt,
		})
	}

	return &types.ListDictTypeResp{
		Total: int(rpcResp.Total),
		List:  list,
	}, nil
}
