// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"

	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/common"
	"go-zero-rpc/sys-rpc/internal/svc"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuTreeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMenuTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuTreeLogic {
	return &GetMenuTreeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetMenuTree 查询全部菜单并构建树形结构（用于菜单管理页）。
//
// 与 GetCurrentUserMenus 的区别：
//   - GetMenuTree：管理员菜单管理界面使用，返回未软删除的全部菜单（不区分可见性、不区分角色）
//   - GetCurrentUserMenus：前端导航菜单使用，按当前用户的角色过滤、过滤隐藏菜单
func (l *GetMenuTreeLogic) GetMenuTree(in *sys.GetMenuTreeReq) (*sys.MenuTreeResp, error) {
	menus, err := l.svcCtx.SysMenuModel.ListAll(l.ctx)
	if err != nil {
		l.Errorf("查询全部菜单失败：%v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	menuTree := common.BuildMenuTree(menus, 0)
	return &sys.MenuTreeResp{List: menuTree}, nil
}
