// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"
	"go-zero-rpc/sys-rpc/pb"

	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

type DeleteDictTypeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDictTypeLogic {
	return &DeleteDictTypeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteDictType 软删除字典类型，并级联软删除该类型下所有字典数据。
func (l *DeleteDictTypeLogic) DeleteDictType(in *pb.DeleteDictTypeReq) (*pb.Empty, error) {
	dictTypeId := in.DictTypeId

	_, err := l.svcCtx.SysDictTypeModel.FindOne(l.ctx, dictTypeId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, xerr.NewCodeError(xerr.ErrNotFound)
		}
		l.Errorf("删除字典类型时查询dictTypeId[%d]是否存在 失败:%v", dictTypeId, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	// 开启事务
	err = l.svcCtx.Orm.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 先删除该类型下所有字典数据
		if err := l.svcCtx.SysDictDataModel.SoftDeleteDictDataByDictTypeIdTrans(l.ctx, tx, dictTypeId); err != nil {
			l.Errorf("删除字典类型[%d]下字典数据失败：%v", dictTypeId, err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		// 2. 再删除该类型
		if err := l.svcCtx.SysDictTypeModel.SoftDeleteDictTypeTrans(l.ctx, tx, dictTypeId); err != nil {
			l.Errorf("删除字典类型失败:%v", err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		return nil
	})
	if err != nil {
		l.Errorf("删除字典类型事务执行失败: %v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	return &pb.Empty{}, nil
}
