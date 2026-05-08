// Code scaffolded by goctl. Safe to edit.
package permissionservicelogic

import (
	"context"

	"go-zero-rpc/common/constants"
	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/svc"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRolesByUserIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRolesByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRolesByUserIdLogic {
	return &GetRolesByUserIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetRolesByUserId 查询用户所拥有的所有启用状态角色编码。
//
// 业务流程：
//  1. 查询用户绑定的角色ID列表
//  2. 逐个查询角色信息，过滤掉已禁用/不存在的角色
//  3. 若包含 admin 角色，直接返回（admin 优先级最高）
func (l *GetRolesByUserIdLogic) GetRolesByUserId(in *sys.GetRolesByUserIdReq) (*sys.GetRolesByUserIdResp, error) {
	if in.UserId == 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	roleIds, err := l.svcCtx.SysUserRoleModel.GetRoleIdsByUserId(l.ctx, in.UserId)
	if err != nil {
		l.Errorf("查询用户[%d]角色失败：%v", in.UserId, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	roleCodes := make([]string, 0, len(roleIds))
	for _, roleId := range roleIds {
		role, rErr := l.svcCtx.SysRoleModel.FindOneByRoleId(l.ctx, roleId)
		if rErr != nil || role == nil || role.Status != 1 {
			continue
		}
		// admin 优先级最高，命中即返回
		if role.Code == constants.RoleCodeAdmin {
			return &sys.GetRolesByUserIdResp{
				RoleCodes: []string{constants.RoleCodeAdmin},
			}, nil
		}
		roleCodes = append(roleCodes, role.Code)
	}

	return &sys.GetRolesByUserIdResp{
		RoleCodes: roleCodes,
	}, nil
}
