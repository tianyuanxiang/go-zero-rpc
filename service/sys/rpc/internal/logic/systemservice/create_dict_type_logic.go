// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"
	"go-zero-rpc/sys-rpc/pb"

	"go-zero-rpc/common/xerr"
	systemmodel "go-zero-rpc/sys-rpc/internal/model"
	"go-zero-rpc/sys-rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CreateDictTypeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateDictTypeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDictTypeLogic {
	return &CreateDictTypeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateDictType 创建字典类型。
//
// 业务流程：
//  1. 校验字典编码唯一性
//  2. 写入字典类型记录
func (l *CreateDictTypeLogic) CreateDictType(in *pb.CreateDictTypeReq) (*pb.CreateDictTypeResp, error) {
	// 1. 检查字典编码唯一性
	exist, err := l.svcCtx.SysDictTypeModel.FindOneByCode(l.ctx, in.DictCode)
	if err != nil && err != sqlx.ErrNotFound {
		l.Errorf("查询字典类型编码[%s]失败：%v", in.DictCode, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}
	if exist != nil {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrDuplicate, "字典类型编码已存在")
	}

	// 2. 插入字典类型
	dictTypeId, err := l.svcCtx.SysDictTypeModel.InsertDictTypeReturningId(l.ctx, &systemmodel.SysDictType{
		Name:   in.DictName,
		Code:   in.DictCode,
		Status: in.Status,
		Remark: in.Remark,
	})
	if err != nil {
		l.Errorf("插入字典类型失败：%v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	return &pb.CreateDictTypeResp{DictTypeId: dictTypeId}, nil
}
