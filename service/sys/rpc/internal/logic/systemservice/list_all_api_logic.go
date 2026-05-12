// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"
	"go-zero-rpc/sys-rpc/pb"

	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAllApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAllApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAllApiLogic {
	return &ListAllApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListAllApi 查询所有未软删除的API接口，用于下拉选择。
func (l *ListAllApiLogic) ListAllApi(in *pb.Empty) (*pb.ListAllApiResp, error) {
	apis, err := l.svcCtx.SysApiModel.ListByIds(l.ctx, nil)
	if err != nil {
		l.Errorf("查询全部接口失败: %v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	list := make([]*pb.ApiOption, 0, len(apis))
	for _, api := range apis {
		list = append(list, &pb.ApiOption{
			Id:      api.Id,
			ApiName: api.ApiName,
			ApiPath: api.ApiPath,
			Remark:  api.Description,
		})
	}
	return &pb.ListAllApiResp{
		ListAll: list,
	}, nil
}
