// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"
	"go-zero-rpc/sys-rpc/pb"

	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/svc"
	casbinpkg "go-zero-rpc/sys-rpc/pkg/casbin"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UpdateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateRole 更新角色基础信息（不含菜单与接口关联）。
//
// 业务流程：
//  1. 检查角色是否存在且未软删除
//  2. 检查变更后的角色编码是否与其他角色冲突
//  3. 拼接动态更新字段并执行更新
//  4. 若编码发生变化，迁移Casbin策略到新编码
func (l *UpdateRoleLogic) UpdateRole(in *pb.UpdateRoleReq) (*pb.Empty, error) {
	// 1. 检查角色是否存在
	existRole, err := l.svcCtx.SysRoleModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, xerr.NewCodeError(xerr.ErrRoleNotFound)
		}
		l.Errorf("查询角色[%d]失败：%v", in.Id, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}
	if existRole.DeletedAt.Valid {
		return nil, xerr.NewCodeError(xerr.ErrRoleNotFound)
	}

	// 2. 如果修改了角色编码，检查新编码是否与其他角色冲突
	if in.HasRoleCode && existRole.Code != in.RoleCode {
		conflictRole, cErr := l.svcCtx.SysRoleModel.FindOneByCode(l.ctx, in.RoleCode)
		if cErr != nil && cErr != sqlx.ErrNotFound {
			l.Errorf("检查角色编码冲突失败：%v", cErr)
			return nil, xerr.NewCodeError(xerr.ErrInternal)
		}
		if conflictRole != nil && conflictRole.Id != in.Id {
			return nil, xerr.NewCodeError(xerr.ErrRoleCodeDuplicate)
		}
	}

	// 3. 拼接更新字段
	updates := make(map[string]interface{})
	if in.HasRoleName {
		updates["name"] = in.RoleName
	}
	if in.HasRoleCode {
		updates["code"] = in.RoleCode
	}
	if in.HasRemark {
		updates["remark"] = in.Remark
	}
	if in.HasSort {
		updates["sort"] = in.Sort
	}

	if len(updates) == 0 {
		return &pb.Empty{}, nil
	}

	if err := l.svcCtx.SysRoleModel.UpdateRoleTrans(l.ctx, in.Id, updates); err != nil {
		l.Errorf("更新角色基本信息失败：%v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	// 4. 如果编码发生变化，迁移Casbin策略到新编码
	if in.HasRoleCode && existRole.Code != in.RoleCode {
		oldPolicies, err := casbinpkg.GetRolePolicies(l.svcCtx.Enforcer, existRole.Code)
		if err != nil {
			l.Errorf("查询角色旧编码[%s]Casbin策略失败：%v", existRole.Code, err)
		} else if len(oldPolicies) > 0 {
			// 提取 [path, method] 部分
			rules := make([][]string, 0, len(oldPolicies))
			for _, p := range oldPolicies {
				if len(p) >= 3 {
					rules = append(rules, []string{p[1], p[2]})
				}
			}
			// 删除旧编码策略
			if err := casbinpkg.RemoveAllPoliciesForRole(l.svcCtx.Enforcer, existRole.Code); err != nil {
				l.Errorf("清除角色旧编码[%s]Casbin策略失败：%v", existRole.Code, err)
			}
			// 写入新编码策略
			if err := casbinpkg.AddRolePolicies(l.svcCtx.Enforcer, in.RoleCode, rules); err != nil {
				l.Errorf("迁移casbin策略到新角色编码[%s]失败: %v", in.RoleCode, err)
			}
		}
	}

	return &pb.Empty{}, nil
}
