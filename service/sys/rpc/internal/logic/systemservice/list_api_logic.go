// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"

	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/svc"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListApiLogic {
	return &ListApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListApi 分页查询API接口列表，支持按分组和关键词筛选。
func (l *ListApiLogic) ListApi(in *sys.ListApiReq) (*sys.ListApiResp, error) {
	apis, count, err := l.svcCtx.SysApiModel.List(l.ctx, int(in.Page), int(in.PageSize), in.Group, in.Keyword)
	if err != nil {
		l.Errorf("查询接口列表失败：%v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	list := make([]*sys.ApiItem, 0, len(apis))
	for _, api := range apis {
		list = append(list, &sys.ApiItem{
			Id:        api.Id,
			ApiName:   api.ApiName,
			ApiPath:   api.ApiPath,
			Method:    api.Method,
			Group:     api.ApiGroup,
			Remark:    api.Description,
			CreatedAt: api.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return &sys.ListApiResp{
		Total: count,
		List:  list,
	}, nil
}
