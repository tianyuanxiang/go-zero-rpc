// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"

	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/svc"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDictDataByTypeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDictDataByTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictDataByTypeLogic {
	return &GetDictDataByTypeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetDictDataByType 根据字典类型ID查询所有字典数据。
func (l *GetDictDataByTypeLogic) GetDictDataByType(in *sys.GetDictDataByTypeReq) (*sys.ListDictDataResp, error) {
	dictData, count, err := l.svcCtx.SysDictDataModel.ListByDictTypeId(l.ctx, in.DictTypeId)
	if err != nil {
		l.Errorf("根据字典类型查询字典数据失败 %v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	list := make([]*sys.DictDataItem, 0, len(dictData))
	for _, item := range dictData {
		list = append(list, &sys.DictDataItem{
			Id:        item.Id,
			DictType:  item.TypeId,
			DictLabel: item.Label,
			DictValue: item.DictValue,
			Sort:      item.Sort,
			Status:    item.Status,
			Remark:    item.Remark,
		})
	}
	return &sys.ListDictDataResp{
		Total: count,
		List:  list,
	}, nil
}
