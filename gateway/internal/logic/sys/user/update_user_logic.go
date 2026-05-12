package user

import (
	"context"

	"go-zero-rpc/common/middleware"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	"go-zero-rpc/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserLogic) UpdateUser(req *types.UpdateUserReq) error {
	userId := middleware.GetUserIdFromCtx(l.ctx)
	if userId == 0 {
		return nil
	}

	rpcReq := &sys.UpdateUserReq{
		Id:         req.Id,
		OperatorId: userId,
		RoleIds:    req.RoleIds,
		HasRoleIds: req.RoleIds != nil,
	}

	rpcReq.HasUsername = req.Username != ""
	rpcReq.Username = req.Username

	if req.Nickname != nil {
		rpcReq.HasNickname = true
		rpcReq.Nickname = *req.Nickname
	}
	if req.Email != nil {
		rpcReq.HasEmail = true
		rpcReq.Email = *req.Email
	}
	if req.Phone != nil {
		rpcReq.HasPhone = true
		rpcReq.Phone = *req.Phone
	}
	if req.Status != nil {
		rpcReq.HasStatus = true
		rpcReq.Status = int64(*req.Status)
	}
	if req.Avatar != nil {
		rpcReq.HasAvatar = true
		rpcReq.Avatar = *req.Avatar
	}
	if req.Remark != nil {
		rpcReq.HasRemark = true
		rpcReq.Remark = *req.Remark
	}

	_, err := l.svcCtx.SysRpc.UpdateUser(l.ctx, rpcReq)
	if err != nil {
		l.Errorf("调用UpdateUser RPC失败, operatorId=%d, targetId=%d, err=%v", userId, req.Id, err)
		return err
	}

	return nil
}
