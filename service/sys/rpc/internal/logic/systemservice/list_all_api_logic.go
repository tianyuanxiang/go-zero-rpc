// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"

	"go-zero-rpc/common/xerr"
	systemmodel "go-zero-rpc/sys-rpc/internal/model"
	"go-zero-rpc/sys-rpc/internal/svc"
	"go-zero-rpc/sys-rpc/pb"

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

func (l *ListAllApiLogic) ListAllApi(in *pb.Empty) (*pb.ListAllApiResp, error) {
	var apis []systemmodel.SysApi
	if err := l.svcCtx.Orm.WithContext(l.ctx).
		Table("sys_api").
		Where("deleted_at IS NULL").
		Order("api_group ASC, api_path ASC, method ASC").
		Find(&apis).Error; err != nil {
		l.Errorf("查询全部接口失败: %v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	const defaultGroup = "未分组"

	groups := make([]*pb.ApiGroupOption, 0)
	groupIndex := make(map[string]int)
	for _, api := range apis {
		group := api.ApiGroup
		if group == "" {
			group = defaultGroup
		}

		idx, ok := groupIndex[group]
		if !ok {
			idx = len(groups)
			groupIndex[group] = idx
			groups = append(groups, &pb.ApiGroupOption{
				Group: group,
				Apis:  make([]*pb.ApiOption, 0),
			})
		}

		groups[idx].Apis = append(groups[idx].Apis, &pb.ApiOption{
			Id:      api.Id,
			ApiName: api.ApiName,
			ApiPath: api.ApiPath,
			Method:  api.Method,
			Remark:  api.Description,
		})
	}

	return &pb.ListAllApiResp{
		Groups: groups,
	}, nil
}
