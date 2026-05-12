// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package log

import (
	"context"

	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearLoginLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClearLoginLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearLoginLogLogic {
	return &ClearLoginLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClearLoginLogLogic) ClearLoginLog() error {
	_, err := l.svcCtx.SysRpc.ClearLoginLog(l.ctx, &sys.ClearLoginLogReq{})
	if err != nil {
		l.Logger.Errorf("调用ClearLoginLog RPC失败, err=%v", err)
		return err
	}

	return nil
}
