// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"

	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/svc"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type GetRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleLogic {
	return &GetRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetRole 根据角色ID查询角色详情。
func (l *GetRoleLogic) GetRole(in *sys.GetRoleReq) (*sys.RoleItem, error) {
	role, err := l.svcCtx.SysRoleModel.FindOneByRoleId(l.ctx, in.RoleId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, xerr.NewCodeError(xerr.ErrRoleNotFound)
		}
		l.Errorf("查询角色[%d]失败：%v", in.RoleId, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}
	return &sys.RoleItem{
		Id:        role.Id,
		RoleCode:  role.Code,
		RoleName:  role.Name,
		Status:    role.Status,
		Sort:      role.Sort,
		Remark:    role.Remark,
		CreatedAt: role.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
