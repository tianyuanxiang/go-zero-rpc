// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package log

import (
	"context"

	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOperLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListOperLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOperLogLogic {
	return &ListOperLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListOperLogLogic) ListOperLog(req *types.ListOperLogReq) (resp *types.ListOperLogResp, err error) {
	rpcReq := &sys.ListOperLogReq{
		Page:         int64(req.Page),
		PageSize:     int64(req.PageSize),
		Keyword:      req.Keyword,
		BusinessType: req.BusinessType,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
	}

	// HTTP 层 Status 是 int 非指针，无法区分零值。这里约定 Status>0 才作为查询条件。
	if req.Status > 0 {
		rpcReq.HasStatus = true
		rpcReq.Status = int64(req.Status)
	}

	rpcResp, err := l.svcCtx.SysRpc.ListOperLog(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("调用ListOperLog RPC失败, err=%v", err)
		return nil, err
	}

	list := make([]types.OperLogItem, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		if item == nil {
			continue
		}
		list = append(list, types.OperLogItem{
			Id:         int(item.Id),
			Title:      item.Title,
			OperType:   item.OperType,
			Method:     item.Method,
			ReqMethod:  item.ReqMethod,
			OperName:   item.OperName,
			DeptName:   item.DeptName,
			ReqUrl:     item.ReqUrl,
			ReqParam:   item.ReqParam,
			RespResult: item.RespResult,
			Status:     int(item.Status),
			Ip:         item.Ip,
			OperTime:   item.OperTime,
		})
	}

	return &types.ListOperLogResp{
		Total: rpcResp.Total,
		List:  list,
	}, nil
}
