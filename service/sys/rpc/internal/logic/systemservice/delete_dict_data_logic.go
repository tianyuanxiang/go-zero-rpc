// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"
	"go-zero-rpc/sys-rpc/pb"

	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DeleteDictDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteDictDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDictDataLogic {
	return &DeleteDictDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteDictData 软删除字典数据。
func (l *DeleteDictDataLogic) DeleteDictData(in *pb.DeleteDictDataReq) (*pb.Empty, error) {
	dictDataId := in.DictDataId

	_, err := l.svcCtx.SysDictDataModel.FindOne(l.ctx, dictDataId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, xerr.NewCodeError(xerr.ErrNotFound)
		}
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	if err := l.svcCtx.SysDictDataModel.SoftDeleteDictDataByDictDataId(l.ctx, dictDataId); err != nil {
		l.Errorf("软删除字典数据[%d]失败：%v", dictDataId, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	return &pb.Empty{}, nil
}
