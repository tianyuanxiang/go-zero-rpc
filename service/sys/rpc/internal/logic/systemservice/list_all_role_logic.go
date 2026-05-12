// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"
	"go-zero-rpc/sys-rpc/pb"

	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAllRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAllRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAllRoleLogic {
	return &ListAllRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListAllRole 查询所有启用状态的角色，用于下拉选择。
func (l *ListAllRoleLogic) ListAllRole(in *pb.Empty) (*pb.ListAllRoleResp, error) {
	roles, err := l.svcCtx.SysRoleModel.ListAll(l.ctx)
	if err != nil {
		l.Errorf("查询全部角色失败: %v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	list := make([]*pb.RoleOption, 0, len(roles))
	for _, role := range roles {
		list = append(list, &pb.RoleOption{
			Id:       role.Id,
			RoleName: role.Name,
			RoleCode: role.Code,
		})
	}

	return &pb.ListAllRoleResp{
		ListAll: list,
	}, nil
}
