// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"

	"go-zero-rpc/common/xerr"
	systemmodel "go-zero-rpc/sys-rpc/internal/model"
	"go-zero-rpc/sys-rpc/internal/svc"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CreateDictDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateDictDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDictDataLogic {
	return &CreateDictDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateDictData 创建字典数据，前置校验所属字典类型存在且未删除。
func (l *CreateDictDataLogic) CreateDictData(in *sys.CreateDictDataReq) (*sys.CreateDictDataResp, error) {
	// 创建数据之前先查是否有该字典类型
	dictType, err := l.svcCtx.SysDictTypeModel.FindOne(l.ctx, in.DictTypeId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			l.Errorf("创建字典数据时，字典类型[%d]不存在", in.DictTypeId)
			return nil, xerr.NewCodeError(xerr.ErrParamInvalid)
		}
		l.Errorf("查询字典类型[%d]失败：%v", in.DictTypeId, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	if dictType.DeletedAt.Valid {
		l.Errorf("创建字典数据时，字典类型[%d]已被删除", in.DictTypeId)
		return nil, xerr.NewCodeError(xerr.ErrNotFound)
	}

	dictDataId, err := l.svcCtx.SysDictDataModel.InsertDictDataReturningId(l.ctx, &systemmodel.SysDictData{
		TypeId:    in.DictTypeId,
		Label:     in.DictLabel,
		DictValue: in.DictValue,
		Sort:      in.Sort,
		Status:    in.Status,
		Remark:    in.Remark,
	})
	if err != nil {
		l.Errorf("插入字典数据失败：%v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	return &sys.CreateDictDataResp{DictDataId: dictDataId}, nil
}
