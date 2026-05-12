package user

import (
	"context"

	"go-zero-rpc/common/middleware"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUserLogic) CreateUser(req *types.CreateUserReq) error {
	userId := middleware.GetUserIdFromCtx(l.ctx)
	if userId == 0 {
		return nil
	}

	_, err := l.svcCtx.SysRpc.CreateUser(l.ctx, &sys.CreateUserReq{
		Username:   req.Username,
		Password:   req.Password,
		Nickname:   req.Nickname,
		Email:      req.Email,
		Phone:      req.Phone,
		Status:     int64(req.Status),
		Avatar:     req.Avatar,
		Remark:     req.Remark,
		RoleIds:    req.RoleIds,
		OperatorId: userId,
	})
	if err != nil {
		l.Logger.Errorf("调用CreateUser RPC失败, operatorId=%d, err=%v", userId, err)
		return err
	}

	return nil
}
