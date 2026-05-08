// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package log

import (
	"context"

	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearOperLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClearOperLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearOperLogLogic {
	return &ClearOperLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClearOperLogLogic) ClearOperLog() error {
	_, err := l.svcCtx.SysRpc.ClearOperLog(l.ctx, &sys.ClearOperLogReq{})
	if err != nil {
		l.Logger.Errorf("调用ClearOperLog RPC失败, err=%v", err)
		return err
	}

	return nil
}
