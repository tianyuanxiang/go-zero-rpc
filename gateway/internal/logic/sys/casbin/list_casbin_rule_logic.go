// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package casbin

import (
	"context"

	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	sys "go-zero-rpc/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCasbinRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListCasbinRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCasbinRuleLogic {
	return &ListCasbinRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListCasbinRuleLogic) ListCasbinRule(req *types.ListCasbinRuleReq) (resp *types.ListCasbinRuleResp, err error) {
	rpcResp, err := l.svcCtx.SysRpc.ListCasbinRule(l.ctx, &sys.ListCasbinRuleReq{
		Page:     int64(req.Page),
		PageSize: int64(req.PageSize),
		RoleCode: req.RoleCode,
		Path:     req.Path,
		Method:   req.Method,
	})
	if err != nil {
		l.Logger.Errorf("Call the ListCasbinRule failed, err=%v", err)
		return nil, err
	}

	list := make([]types.CasbinRuleItem, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		if item == nil {
			continue
		}
		list = append(list, types.CasbinRuleItem{
			Id:       item.Id,
			Ptype:    item.Ptype,
			RoleCode: item.RoleCode,
			ApiPath:  item.ApiPath,
			Method:   item.Method,
		})
	}

	return &types.ListCasbinRuleResp{
		Total: rpcResp.Total,
		List:  list,
	}, nil
}
