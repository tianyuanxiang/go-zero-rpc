package user

import (
	"context"

	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserLogic {
	return &ListUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUserLogic) ListUser(req *types.ListUserReq) (resp *types.ListUserResp, err error) {
	rpcReq := &sys.ListUserReq{
		Page:     int64(req.Page),
		PageSize: int64(req.PageSize),
		Keyword:  req.Keyword,
	}
	if req.Status != nil {
		rpcReq.HasStatus = true
		rpcReq.Status = int64(*req.Status)
	}

	rpcResp, err := l.svcCtx.SysRpc.ListUser(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("调用ListUser RPC失败, err=%v", err)
		return nil, err
	}

	list := make([]types.UserItem, 0, len(rpcResp.List))
	for _, user := range rpcResp.List {
		list = append(list, types.UserItem{
			Id:        user.Id,
			Username:  user.Username,
			Nickname:  user.Nickname,
			Email:     user.Email,
			Phone:     user.Phone,
			Status:    int(user.Status),
			Avatar:    user.Avatar,
			Remark:    user.Remark,
			Roles:     user.Roles,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return &types.ListUserResp{
		Total: rpcResp.Total,
		List:  list,
	}, nil
}
