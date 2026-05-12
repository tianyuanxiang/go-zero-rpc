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

type CreateApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateApiLogic {
	return &CreateApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateApi 创建新API接口。
//
// 业务流程：
//  1. 校验路径+方法组合唯一性
//  2. 写入接口数据
func (l *CreateApiLogic) CreateApi(in *pb.CreateApiReq) (*pb.CreateApiResp, error) {
	// 1. 检查路径+方法组合唯一性
	exist, err := l.svcCtx.SysApiModel.FindOneByApiPathMethod(l.ctx, in.ApiPath, in.Method)
	if err != nil && err != sqlx.ErrNotFound {
		l.Errorf("查询接口[%s %s]是否存在失败：%v", in.Method, in.ApiPath, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}
	if exist != nil {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrDuplicate, "该路径和方法组合已存在")
	}

	// 2. 插入接口
	apiId, err := l.svcCtx.SysApiModel.InsertApiReturningId(l.ctx, &systemmodel.SysApi{
		ApiPath:     in.ApiPath,
		ApiName:     in.ApiName,
		Method:      in.Method,
		ApiGroup:    in.Group,
		Description: in.Remark,
	})
	if err != nil {
		l.Errorf("插入接口记录失败：%v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	return &pb.CreateApiResp{ApiId: apiId}, nil
}
