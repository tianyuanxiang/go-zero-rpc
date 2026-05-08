// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"

	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/svc"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRoleLogic {
	return &ListRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListRole 分页查询角色列表，支持关键词模糊检索。
func (l *ListRoleLogic) ListRole(in *sys.ListRoleReq) (*sys.ListRoleResp, error) {
	roles, count, err := l.svcCtx.SysRoleModel.List(l.ctx, int(in.Page), int(in.PageSize), in.Keyword)
	if err != nil {
		l.Errorf("查询角色列表失败 %v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	list := make([]*sys.RoleItem, 0, len(roles))
	for _, role := range roles {
		list = append(list, &sys.RoleItem{
			Id:        role.Id,
			RoleName:  role.Name,
			RoleCode:  role.Code,
			Status:    role.Status,
			Sort:      role.Sort,
			Remark:    role.Remark,
			CreatedAt: role.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &sys.ListRoleResp{
		Total: count,
		List:  list,
	}, nil
}
