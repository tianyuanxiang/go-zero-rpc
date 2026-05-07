package user

import (
	"context"

	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLogic) GetUser(req *types.GetUserReq) (resp *types.UserItem, err error) {
	user, err := l.svcCtx.SysRpc.GetUser(l.ctx, &sys.GetUserReq{
		UserId: req.Id,
	})
	if err != nil {
		l.Logger.Errorf("调用GetUser RPC失败, userId=%d, err=%v", req.Id, err)
		return nil, err
	}

	return &types.UserItem{
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
	}, nil
}
