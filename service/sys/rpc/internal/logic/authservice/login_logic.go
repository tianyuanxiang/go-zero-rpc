package authservicelogic

import (
	"context"
	"go-zero-rpc/common/constants"
	"go-zero-rpc/common/jwtx"
	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/common"
	sysmodel "go-zero-rpc/sys-rpc/internal/model"
	"go-zero-rpc/sys-rpc/pkg/encrypt"
	"time"

	"go-zero-rpc/sys-rpc/internal/svc"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Login 处理用户登录
//
// 业务流程与单体版完全一致：
// 1. 按用户名查询用户（不存在 → 错误）
// 2. 校验账号状态（禁用 → 错误）
// 3. bcrypt 校验密码
// 4. 生成 JWT 双 Token
// 5. 查询用户角色 + 菜单
// 6. 异步记录登录日志
// 7. 返回 LoginResp

func (l *LoginLogic) Login(in *sys.LoginReq) (*sys.LoginResp, error) {
	// 1. 查用户
	user, err := l.svcCtx.SysUserModel.FindOneByUsername(l.ctx, in.Username)
	if err != nil {
		if err == sysmodel.ErrNotFound {
			l.recordLoginLog(in.Username, 0, in.ClientIp, in.UserAgent, 0, "用户名或密码错误")
			return nil, xerr.NewCodeError(xerr.ErrPasswordWrong)
		}
		l.Errorf("查询用户失败: %v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	// 2. 状态
	if user.Status != 1 {
		l.recordLoginLog(in.Username, user.Id, in.ClientIp, in.UserAgent, 0, "账号已被禁用")
		return nil, xerr.NewCodeError(xerr.ErrAccountDisabled)
	}

	// 3. 密码
	if !encrypt.CheckPassword(in.Password, user.Password) {
		l.recordLoginLog(in.Username, user.Id, in.ClientIp, in.UserAgent, 0, "密码错误")
		return nil, xerr.NewCodeError(xerr.ErrPasswordWrong)
	}

	// 4. Token
	accessToken, err := jwtx.GenerateToken(user.Id, user.Username,
		l.svcCtx.Config.JwtAuth.AccessSecret, l.svcCtx.Config.JwtAuth.AccessExpire)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}
	refreshToken, err := jwtx.GenerateRefreshToken(user.Id, user.Username,
		l.svcCtx.Config.JwtAuth.AccessSecret, l.svcCtx.Config.JwtAuth.RefreshExpire)
	if err != nil {
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	// 5. 查询角色ID
	roleIds, err := l.svcCtx.SysUserRoleModel.GetRoleIdsByUserId(l.ctx, user.Id)
	if err != nil {
		l.Logger.Errorf("查询用户[%s]的角色信息失败：%v", user.Username, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	roleCodes := make([]string, 0, len(roleIds)+1)
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

	// 6. 查询用户的菜单权限（合并所有角色的菜单）
	var menus []*sysmodel.SysMenu // 创建切片
	if isAdmin {
		// 直接查询所有菜单
		menus, err = l.svcCtx.SysMenuModel.ListAll(l.ctx)
		if err != nil {
			l.Logger.Errorf("查询超管所有菜单失败：%v", err)
			menus = []*sysmodel.SysMenu{}
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
			menus = []*sysmodel.SysMenu{}
		}
	}

	// 将菜单列表构建为树形结构
	menuTree := common.BuildMenuTree(menus, 0)

	// 6. 登录日志
	l.recordLoginLog(in.Username, user.Id, in.ClientIp, in.UserAgent, 1, "登录成功")

	// 7. 组装返回（types.UserInfo 改成 sys.UserInfo）
	return &sys.LoginResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    l.svcCtx.Config.JwtAuth.AccessExpire,
		UserInfo: &sys.UserInfo{
			UserId:   user.Id,
			Username: user.Username,
			Email:    user.Email,
			Phone:    user.Phone,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Roles:    roleCodes,
			Menus:    menuTree, // []*sys.MenuItem
		},
	}, nil
}

func (l *LoginLogic) recordLoginLog(username string, userId int64, ip, ua string, status int64, msg string) {
	go func() {
		// 浏览器/操作系统识别逻辑直接复制原项目 internal/common/browser.go
		// ...
		loginLog := &sysmodel.SysLoginLog{
			UserId:    userId,
			Username:  username,
			Ip:        ip,
			Browser:   "", // common.ExtractBrowser(ua)
			Os:        "", // common.ExtractOS(ua)
			Status:    status,
			Msg:       msg,
			LoginTime: time.Now(),
		}
		l.svcCtx.SysLoginLogModel.Insert(context.Background(), loginLog)
	}()
}
