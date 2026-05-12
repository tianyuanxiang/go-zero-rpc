package permissionservicelogic

import (
	"context"
	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/pb"

	"go-zero-rpc/sys-rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsTokenRevokedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsTokenRevokedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsTokenRevokedLogic {
	return &IsTokenRevokedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsTokenRevokedLogic) IsTokenRevoked(in *pb.IsBlackListReq) (*pb.IsBlackListResp, error) {

	exists, err := l.svcCtx.RDB.Exists(l.ctx, in.BlacklistKey).Result()
	if err != nil {
		l.Logger.Errorf("Query to the redis blacklist failed, err:%v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}
	if exists > 0 {
		l.Logger.Infof("Token has expired.")
		return &pb.IsBlackListResp{
			IsBlack: true,
		}, nil
	}
	return &pb.IsBlackListResp{
		IsBlack: false,
	}, nil
}
