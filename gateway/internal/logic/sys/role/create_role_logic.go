// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role

import (
	"context"

	"go-zero-rpc/common/middleware"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRoleLogic) CreateRole(req *types.CreateRoleReq) error {
	operatorId := middleware.GetUserIdFromCtx(l.ctx)

	_, err := l.svcCtx.SysRpc.CreateRole(l.ctx, &sys.CreateRoleReq{
		RoleName:   req.RoleName,
		RoleCode:   req.RoleCode,
		Sort:       req.Sort,
		Remark:     req.Remark,
		MenuIds:    req.MenuIds,
		ApiIds:     req.ApiIds,
		OperatorId: operatorId,
	})
	if err != nil {
		l.Logger.Errorf("调用CreateRole RPC失败, operatorId=%d, err=%v", operatorId, err)
		return err
	}

	return nil
}
