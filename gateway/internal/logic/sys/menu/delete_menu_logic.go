// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package menu

import (
	"context"

	"go-zero-rpc/common/middleware"
	"go-zero-rpc/gateway/internal/svc"
	sys "go-zero-rpc/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuLogic {
	return &DeleteMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMenuLogic) DeleteMenu(menuId int64) error {
	userId := middleware.GetUserIdFromCtx(l.ctx)

	_, err := l.svcCtx.SysRpc.DeleteMenu(l.ctx, &sys.DeleteMenuReq{
		MenuId:     menuId,
		OperatorId: userId,
	})
	if err != nil {
		l.Logger.Errorf("Call the DeleteMenu failed, operatorId=%d, menuId=%d, err=%v", userId, menuId, err)
		return err
	}

	return nil
}
