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

type UpdateDictDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateDictDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDictDataLogic {
	return &UpdateDictDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateDictData 更新字典数据。
func (l *UpdateDictDataLogic) UpdateDictData(in *pb.UpdateDictDataReq) (*pb.Empty, error) {
	dictData, err := l.svcCtx.SysDictDataModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, xerr.NewCodeError(xerr.ErrNotFound)
		}
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	if dictData.DeletedAt.Valid {
		l.Errorf("字典数据[%d]已删除", in.Id)
		return nil, xerr.NewCodeError(xerr.ErrParamInvalid)
	}

	// 拼接动态更新字段
	updates := make(map[string]interface{})
	if in.DictTypeId != 0 {
		updates["type_id"] = in.DictTypeId
	}
	if in.DictLabel != "" {
		updates["label"] = in.DictLabel
	}
	if in.DictValue != "" {
		updates["dict_value"] = in.DictValue
	}
	if in.HasSort {
		updates["sort"] = in.Sort
	}
	if in.HasStatus {
		updates["status"] = in.Status
	}
	if in.HasRemark {
		updates["remark"] = in.Remark
	}

	if len(updates) == 0 {
		l.Error("更新字段为空")
		return nil, xerr.NewCodeError(xerr.ErrParamInvalid)
	}

	if err := l.svcCtx.SysDictDataModel.UpdateDictData(l.ctx, in.Id, updates); err != nil {
		l.Errorf("更新DictData[%d]失败：%v", in.Id, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}
	return &pb.Empty{}, nil
}
