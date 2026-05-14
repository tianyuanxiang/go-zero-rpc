// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package log

import (
	"context"

	"go-zero-rpc/gateway/internal/svc"
	sys "go-zero-rpc/sys-rpc/pb"

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
		l.Logger.Errorf("Call the ClearOperLog failed, err=%v", err)
		return err
	}

	return nil
}
