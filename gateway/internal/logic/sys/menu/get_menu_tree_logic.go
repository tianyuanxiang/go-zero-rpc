// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package menu

import (
	"context"

	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	sys "go-zero-rpc/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMenuTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuTreeLogic {
	return &GetMenuTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMenuTreeLogic) GetMenuTree() (resp *types.MenuTreeResp, err error) {
	rpcResp, err := l.svcCtx.SysRpc.GetMenuTree(l.ctx, &sys.GetMenuTreeReq{})
	if err != nil {
		l.Logger.Errorf("Call the GetMenuTree failed, err=%v", err)
		return nil, err
	}

	return &types.MenuTreeResp{
		List: convertMenuItems(rpcResp.List),
	}, nil
}

// convertMenuItems maps RPC menu nodes to HTTP response nodes.
func convertMenuItems(in []*sys.MenuItem) []types.MenuItem {
	out := make([]types.MenuItem, 0, len(in))
	for _, item := range in {
		if item == nil {
			continue
		}
		out = append(out, types.MenuItem{
			Id:        item.Id,
			ParentId:  item.ParentId,
			MenuName:  item.MenuName,
			MenuType:  item.MenuType,
			Path:      item.Path,
			Component: item.Component,
			Icon:      item.Icon,
			Sort:      item.Sort,
			Perms:     item.Perms,
			Status:    int(item.Status),
			Children:  convertMenuItems(item.Children),
		})
	}
	return out
}
