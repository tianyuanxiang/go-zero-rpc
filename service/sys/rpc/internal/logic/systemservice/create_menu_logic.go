// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"
	"go-zero-rpc/sys-rpc/pb"

	"go-zero-rpc/common/xerr"
	systemmodel "go-zero-rpc/sys-rpc/internal/model"
	"go-zero-rpc/sys-rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMenuLogic {
	return &CreateMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateMenu 创建新菜单。
//
// 业务流程：
//  1. 若指定父菜单ID，校验父菜单存在且不是按钮类型
//  2. 写入菜单数据
func (l *CreateMenuLogic) CreateMenu(in *pb.CreateMenuReq) (*pb.CreateMenuResp, error) {
	// 1. 确认父菜单是否存在
	if in.ParentId > 0 {
		menus, err := l.svcCtx.SysMenuModel.ListByIds(l.ctx, []int64{in.ParentId})
		if err != nil {
			l.Errorf("查询父菜单[%d]失败：%v", in.ParentId, err)
			return nil, xerr.NewCodeError(xerr.ErrInternal)
		}
		if len(menus) == 0 {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrMenuNotFound, "父菜单不存在或已被删除")
		}
		// 额外校验：父菜单不能是按钮类型
		if menus[0].MenuType == 2 {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "按钮类型的菜单不能作为父菜单")
		}
	}

	// 2. 插入菜单
	menuId, err := l.svcCtx.SysMenuModel.InsertMenuReturningId(l.ctx, &systemmodel.SysMenu{
		ParentId:   in.ParentId,
		Name:       in.MenuName,
		MenuType:   in.MenuType,
		MenuPath:   in.Path,
		Component:  in.Component,
		Icon:       in.Icon,
		Sort:       in.Sort,
		Permission: in.Perms,
		Status:     in.Status,
		Visible:    in.Visible,
		Remark:     in.Remark,
	})
	if err != nil {
		l.Errorf("插入菜单记录失败：%v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	return &pb.CreateMenuResp{MenuId: menuId}, nil
}
