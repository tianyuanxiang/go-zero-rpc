// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"
	"go-zero-rpc/sys-rpc/pb"

	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/common"
	"go-zero-rpc/sys-rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UpdateApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateApiLogic {
	return &UpdateApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateApi 更新API接口信息。
//
// 业务流程：
//  1. 校验接口存在性
//  2. 拼接动态更新字段并执行更新
//  3. 若 path 或 method 发生变化，重建关联角色的Casbin策略
func (l *UpdateApiLogic) UpdateApi(in *pb.UpdateApiReq) (*pb.Empty, error) {
	oldApi, err := l.svcCtx.SysApiModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, xerr.NewCodeError(xerr.ErrMenuNotFound)
		}
		l.Errorf("更新api接口查询旧记录失败：%v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}
	if oldApi.DeletedAt.Valid {
		return nil, xerr.NewCodeError(xerr.ErrNotFound)
	}

	// 拼接动态更新字段
	updates := make(map[string]interface{})
	if in.HasApiPath {
		updates["api_path"] = in.ApiPath
	}
	if in.HasMethod {
		updates["method"] = in.Method
	}
	if in.HasGroup {
		updates["api_group"] = in.Group
	}
	if in.HasRemark {
		updates["description"] = in.Remark
	}
	if in.HasApiName {
		updates["api_name"] = in.ApiName
	}

	if len(updates) == 0 {
		l.Error("更新字段为空")
		return nil, xerr.NewCodeError(xerr.ErrParamInvalid)
	}

	if err := l.svcCtx.SysApiModel.UpdateApi(l.ctx, in.Id, updates); err != nil {
		l.Errorf("更新api[%d]失败：%v", in.Id, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	// 判断path或method是否变更，若变更则同步Casbin
	pathChanged := in.HasApiPath && oldApi.ApiPath != in.ApiPath
	methodChanged := in.HasMethod && oldApi.Method != in.Method
	if pathChanged || methodChanged {
		roleIds, err := l.svcCtx.SysRoleApiModel.ListRoleIdsByApiId(l.ctx, in.Id)
		if err != nil {
			l.Errorf("查询API[%d]关联角色失败：%v", in.Id, err)
		} else {
			common.RebuildCasbinByRoleIds(
				l.ctx,
				l.svcCtx.Enforcer,
				l.svcCtx.SysRoleModel,
				l.svcCtx.SysRoleApiModel,
				l.svcCtx.SysApiModel,
				roleIds,
			)
		}
	}

	return &pb.Empty{}, nil
}
