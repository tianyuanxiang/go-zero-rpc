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

type UpdateDictTypeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDictTypeLogic {
	return &UpdateDictTypeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateDictType 更新字典类型。
func (l *UpdateDictTypeLogic) UpdateDictType(in *pb.UpdateDictTypeReq) (*pb.Empty, error) {
	dictType, err := l.svcCtx.SysDictTypeModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, xerr.NewCodeError(xerr.ErrNotFound)
		}
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	if dictType.DeletedAt.Valid {
		l.Errorf("字典类型[%d]已删除", in.Id)
		return nil, xerr.NewCodeError(xerr.ErrParamInvalid)
	}

	// 拼接动态更新字段
	updates := make(map[string]interface{})
	if in.HasDictCode {
		updates["code"] = in.DictCode
	}
	if in.HasDictName {
		updates["name"] = in.DictName
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

	if err := l.svcCtx.SysDictTypeModel.UpdateDictType(l.ctx, in.Id, updates); err != nil {
		l.Errorf("更新DictType[%d]失败：%v", in.Id, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}
	return &pb.Empty{}, nil
}
