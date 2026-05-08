// Code scaffolded by goctl. Safe to edit.
package permissionservicelogic

import (
	"context"
	"fmt"

	"go-zero-rpc/common/constants"
	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/svc"
	casbinpkg "go-zero-rpc/sys-rpc/pkg/casbin"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckPermissionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckPermissionLogic {
	return &CheckPermissionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CheckPermission 校验指定用户是否对某 path+method 有访问权限。
//
// 业务流程：
//  1. 查询用户所有角色ID
//  2. 遍历角色取出 role.Code
//     - 若任一角色为 admin，直接放行
//     - 否则用 casbin 逐个校验 (path, method)，任一通过即放行
//
// 参数：
//   - in : 校验请求体（user_id / path / method）
//
// 返回：
//   - *sys.CheckPermissionResp : allowed=true 表示有权限
//   - error                    : 业务错误
func (l *CheckPermissionLogic) CheckPermission(in *sys.CheckPermissionReq) (*sys.CheckPermissionResp, error) {
	if in.UserId == 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	l.Logger.Infof("\033[31m%s\033[0m", "***************************** casbinMiddleware start **********************")
	// 1. 查询用户的所有角色ID
	roleIds, err := l.svcCtx.SysUserRoleModel.GetRoleIdsByUserId(l.ctx, in.UserId)
	if err != nil {
		l.Errorf("查询用户[%d]角色失败：%v", in.UserId, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	if len(roleIds) == 0 {
		l.Logger.Infof("\u001B[31m%s\u001B[0m", "********************************** casbinMiddleware end *********************************")
		return &sys.CheckPermissionResp{Allowed: false, Reason: "roleIds is not found."}, nil
	}

	// 2. 校验权限
	var roleCodes []string
	for _, roleId := range roleIds {
		role, rErr := l.svcCtx.SysRoleModel.FindOneByRoleId(l.ctx, roleId)
		if rErr != nil || role == nil || role.Status != 1 {
			continue
		}
		// 2.1 admin 角色直接放行
		if role.Code == constants.RoleCodeAdmin {
			l.Logger.Infof("\u001B[31m%s\u001B[0m", "********************************** casbinMiddleware end *********************************")
			return &sys.CheckPermissionResp{Allowed: true}, nil
		}

		// 2.2 普通角色走 casbin 校验
		ok, cErr := casbinpkg.CheckPermission(l.svcCtx.Enforcer, role.Code, in.Path, in.Method)
		if cErr != nil {
			l.Errorf("Casbin 校验角色[%s]访问[%s %s]失败：%v", role.Code, in.Method, in.Path, cErr)
			continue
		}
		if ok {
			l.Logger.Infof("\u001B[31m%s\u001B[0m", "********************************** casbinMiddleware end *********************************")
			return &sys.CheckPermissionResp{Allowed: true}, nil
		}
		roleCodes = append(roleCodes, role.Code)
	}

	l.Logger.Infof("\u001B[31m%s\u001B[0m", "********************************** casbinMiddleware end *********************************")
	return &sys.CheckPermissionResp{Allowed: false,
		Reason: fmt.Sprintf("Code: %s, Path: %s, Method: %s is denied.", roleCodes, in.Path, in.Method)}, nil
}
