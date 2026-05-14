// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package menu

import (
	"context"

	"go-zero-rpc/common/middleware"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	sys "go-zero-rpc/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuLogic {
	return &UpdateMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMenuLogic) UpdateMenu(req *types.UpdateMenuReq) error {
	operatorId := middleware.GetUserIdFromCtx(l.ctx)

	rpcReq := &sys.UpdateMenuReq{
		Id:         req.Id,
		OperatorId: operatorId,
	}

	if req.ParentId != nil {
		rpcReq.HasParentId = true
		rpcReq.ParentId = *req.ParentId
	}
	if req.MenuName != nil {
		rpcReq.HasMenuName = true
		rpcReq.MenuName = *req.MenuName
	}
	if req.MenuType != nil {
		rpcReq.HasMenuType = true
		rpcReq.MenuType = int64(*req.MenuType)
	}
	if req.Path != nil {
		rpcReq.HasPath = true
		rpcReq.Path = *req.Path
	}
	if req.Component != nil {
		rpcReq.HasComponent = true
		rpcReq.Component = *req.Component
	}
	if req.Icon != nil {
		rpcReq.HasIcon = true
		rpcReq.Icon = *req.Icon
	}
	if req.Sort != nil {
		rpcReq.HasSort = true
		rpcReq.Sort = int64(*req.Sort)
	}
	if req.Perms != nil {
		rpcReq.HasPerms = true
		rpcReq.Perms = *req.Perms
	}
	if req.Status != nil {
		rpcReq.HasStatus = true
		rpcReq.Status = int64(*req.Status)
	}
	if req.Visible != nil {
		rpcReq.HasVisible = true
		rpcReq.Visible = int64(*req.Visible)
	}
	if req.Remark != nil {
		rpcReq.HasRemark = true
		rpcReq.Remark = *req.Remark
	}

	_, err := l.svcCtx.SysRpc.UpdateMenu(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Call the UpdateMenu failed, operatorId=%d, menuId=%d, err=%v", operatorId, req.Id, err)
		return err
	}

	return nil
}
