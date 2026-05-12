// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"
	"go-zero-rpc/sys-rpc/pb"

	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDictTypeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDictTypeLogic {
	return &ListDictTypeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListDictType 分页查询字典类型列表，支持关键词模糊检索。
func (l *ListDictTypeLogic) ListDictType(in *pb.ListDictTypeReq) (*pb.ListDictTypeResp, error) {
	dictTypes, count, err := l.svcCtx.SysDictTypeModel.List(l.ctx, int(in.Page), int(in.PageSize), in.Keyword)
	if err != nil {
		l.Errorf("查询字典类型列表失败：%v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	list := make([]*pb.DictTypeItem, 0, len(dictTypes))
	for _, dictType := range dictTypes {
		list = append(list, &pb.DictTypeItem{
			Id:        dictType.Id,
			DictName:  dictType.Name,
			DictCode:  dictType.Code,
			Remark:    dictType.Remark,
			Status:    dictType.Status,
			CreatedAt: dictType.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return &pb.ListDictTypeResp{
		Total: count,
		List:  list,
	}, nil
}
