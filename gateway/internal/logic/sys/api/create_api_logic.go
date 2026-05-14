// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package api

import (
	"context"

	"go-zero-rpc/common/middleware"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	sys "go-zero-rpc/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateApiLogic {
	return &CreateApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateApiLogic) CreateApi(req *types.CreateApiReq) error {
	operatorId := middleware.GetUserIdFromCtx(l.ctx)

	_, err := l.svcCtx.SysRpc.CreateApi(l.ctx, &sys.CreateApiReq{
		ApiName:    req.ApiName,
		ApiPath:    req.ApiPath,
		Method:     req.Method,
		Group:      req.Group,
		Remark:     req.Remark,
		OperatorId: operatorId,
	})
	if err != nil {
		l.Logger.Errorf("Call the CreateApi failed, operatorId=%d, err=%v", operatorId, err)
		return err
	}

	return nil
}
