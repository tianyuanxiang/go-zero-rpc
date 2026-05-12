// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package menu

import (
	"context"

	"go-zero-rpc/common/middleware"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	"go-zero-rpc/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMenuLogic {
	return &CreateMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateMenuLogic) CreateMenu(req *types.CreateMenuReq) error {
	operatorId := middleware.GetUserIdFromCtx(l.ctx)

	_, err := l.svcCtx.SysRpc.CreateMenu(l.ctx, &sys.CreateMenuReq{
		ParentId:   req.ParentId,
		MenuName:   req.MenuName,
		MenuType:   req.MenuType,
		Path:       req.Path,
		Component:  req.Component,
		Icon:       req.Icon,
		Sort:       int64(req.Sort),
		Perms:      req.Perms,
		Status:     int64(req.Status),
		Remark:     req.Remark,
		OperatorId: operatorId,
	})
	if err != nil {
		l.Logger.Errorf("调用CreateMenu RPC失败, operatorId=%d, err=%v", operatorId, err)
		return err
	}

	return nil
}
