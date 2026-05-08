// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"

	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/common"
	"go-zero-rpc/sys-rpc/internal/svc"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

type DeleteApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteApiLogic {
	return &DeleteApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteApi 软删除API接口，并级联清理 role_api 关联以及Casbin策略。
//
// 业务流程：
//  1. 校验接口存在性
//  2. 提前查询接口绑定的所有角色（用于事务后重建Casbin策略）
//  3. 事务内：软删除接口、删除 role_api 关联
//  4. 事务成功后重建 Casbin 策略
func (l *DeleteApiLogic) DeleteApi(in *sys.DeleteApiReq) (*sys.Empty, error) {
	apiId := in.ApiId

	oldApi, err := l.svcCtx.SysApiModel.FindOne(l.ctx, apiId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, xerr.NewCodeError(xerr.ErrNotFound)
		}
		l.Errorf("删除api时查询apiId[%d]是否存在失败:%v", apiId, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}
	if oldApi.DeletedAt.Valid {
		return nil, xerr.NewCodeError(xerr.ErrMenuNotFound)
	}

	roleIds, err := l.svcCtx.SysRoleApiModel.ListRoleIdsByApiId(l.ctx, apiId)
	if err != nil {
		l.Errorf("查询API[%d]关联角色失败：%v", apiId, err)
	}

	// 1. 开启删除事务
	err = l.svcCtx.Orm.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		// 1.1 删除api
		if err := l.svcCtx.SysApiModel.SoftDeleteApiTrans(l.ctx, tx, apiId); err != nil {
			l.Errorf("删除api失败：%v", err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		// 1.2 删除 role_api 关联关系
		if err := l.svcCtx.SysRoleApiModel.DeleteRoleApiByApiIdTrans(l.ctx, tx, apiId); err != nil {
			l.Errorf("删除role_api关联关系失败%v", err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		return nil
	})
	if err != nil {
		l.Errorf("创建删除api事务执行失败: %v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	// 2. 清除Casbin关联（事务成功后执行）
	common.RebuildCasbinByRoleIds(
		l.ctx,
		l.svcCtx.Enforcer,
		l.svcCtx.SysRoleModel,
		l.svcCtx.SysRoleApiModel,
		l.svcCtx.SysApiModel,
		roleIds,
	)

	return &sys.Empty{}, nil
}
