// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package log

import (
	"context"

	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	"go-zero-rpc/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLoginLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLoginLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLoginLogLogic {
	return &ListLoginLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLoginLogLogic) ListLoginLog(req *types.ListLoginLogReq) (resp *types.ListLoginLogResp, err error) {
	rpcReq := &sys.ListLoginLogReq{
		Page:      int64(req.Page),
		PageSize:  int64(req.PageSize),
		Keyword:   req.Keyword,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	// HTTP 层 Status 是 int 非指针，无法区分零值。这里约定 Status>0 才作为查询条件。
	if req.Status > 0 {
		rpcReq.HasStatus = true
		rpcReq.Status = int64(req.Status)
	}

	rpcResp, err := l.svcCtx.SysRpc.ListLoginLog(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("调用ListLoginLog RPC失败, err=%v", err)
		return nil, err
	}

	list := make([]types.LoginLogItem, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		if item == nil {
			continue
		}
		list = append(list, types.LoginLogItem{
			Id:        int(item.Id),
			UserId:    int(item.UserId),
			Username:  item.Username,
			Ip:        item.Ip,
			Location:  item.Location,
			UserAgent: item.UserAgent,
			OS:        item.Os,
			Status:    int(item.Status),
			Msg:       item.Msg,
			LoginTime: item.LoginTime,
		})
	}

	return &types.ListLoginLogResp{
		Total: int(rpcResp.Total),
		List:  list,
	}, nil
}
