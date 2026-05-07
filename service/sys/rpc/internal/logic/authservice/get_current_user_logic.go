package authservicelogic

import (
	"context"
	"go-zero-rpc/common/constants"
	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/common"
	systemmodel "go-zero-rpc/sys-rpc/internal/model"

	"go-zero-rpc/sys-rpc/internal/svc"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type GetCurrentUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCurrentUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrentUserLogic {
	return &GetCurrentUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCurrentUserLogic) GetCurrentUser(in *sys.GetCurrentUserReq) (*sys.GetCurrentUserResp, error) {
	userId := in.UserId
	if userId == 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	// 查询用户基本信息
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, userId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, xerr.NewCodeError(xerr.ErrUserNotFound)
		}
		l.Errorf("查询用户[%d]失败：%v", userId, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	// 查询用户角色
	roleIds, err := l.svcCtx.SysUserRoleModel.GetRoleIdsByUserId(l.ctx, userId)
	if err != nil {
		l.Errorf("查询用户[%d]角色失败：%v", userId, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	roleCodes := make([]string, 0, len(roleIds))
	var isAdmin bool

	for _, roleId := range roleIds {
		role, roleErr := l.svcCtx.SysRoleModel.FindOneByRoleId(l.ctx, roleId)
		if roleErr != nil || role == nil || role.Status != 1 {
			continue
		}
		if role.Code == constants.RoleCodeAdmin {
			isAdmin = true
			roleCodes = append(roleCodes, role.Code)
			break
		}
		roleCodes = append(roleCodes, role.Code)
	}

	var menus []*systemmodel.SysMenu // 创建切片
	if isAdmin {
		// 直接查询所有菜单
		menus, err = l.svcCtx.SysMenuModel.ListAll(l.ctx)
		if err != nil {
			l.Logger.Errorf("查询超管所有菜单失败：%v", err)
			menus = []*systemmodel.SysMenu{}
		}
	} else {
		menuIds, err := l.svcCtx.SysRoleMenuModel.GetMenuIdsByRoleIds(l.ctx, roleIds)
		if err != nil {
			l.Logger.Errorf("查询用户菜单失败：%v", err)
			menuIds = []int64{}
		}
		menus, err = l.svcCtx.SysMenuModel.ListByIds(l.ctx, menuIds)
		if err != nil {
			l.Logger.Errorf("查询菜单列表失败：%v", err)
			menus = []*systemmodel.SysMenu{}
		}
	}

	menuTree := common.BuildMenuTree(menus, 0)

	return &sys.GetCurrentUserResp{
		UserInfo: &sys.UserInfo{
			UserId:   user.Id,
			Username: user.Username,
			Nickname: user.Nickname,
			Email:    user.Email,
			Phone:    user.Phone,
			Avatar:   user.Avatar,
			Roles:    roleCodes,
			Menus:    menuTree,
		},
	}, nil
}
