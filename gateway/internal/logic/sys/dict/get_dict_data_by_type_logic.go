// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package dict

import (
	"context"

	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDictDataByTypeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDictDataByTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictDataByTypeLogic {
	return &GetDictDataByTypeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDictDataByTypeLogic) GetDictDataByType(dictTypeId int64) (resp *types.ListDictDataResp, err error) {
	result, err := l.svcCtx.SysRpc.GetDictDataByType(l.ctx, &sys.GetDictDataByTypeReq{
		DictTypeId: dictTypeId,
	})
	if err != nil {
		l.Logger.Errorf("调用GetDictDataByType RPC失败, dictTypeId=%d, err=%v", dictTypeId, err)
		return nil, err
	}

	list := make([]types.DictDataItem, 0, len(result.List))
	for _, item := range result.List {
		list = append(list, types.DictDataItem{
			Id:        int(item.Id),
			DictType:  int(item.DictType),
			DictLabel: item.DictLabel,
			DictValue: item.DictValue,
			Sort:      int(item.Sort),
			Remark:    item.Remark,
			Status:    int(item.Status),
		})
	}

	return &types.ListDictDataResp{
		Total: result.Total,
		List:  list,
	}, nil
}
