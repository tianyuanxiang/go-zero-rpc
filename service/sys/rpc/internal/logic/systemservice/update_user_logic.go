package systemservicelogic

import (
	"context"
	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/pb"

	"go-zero-rpc/sys-rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserLogic) UpdateUser(in *pb.UpdateUserReq) (*pb.Empty, error) {

	// 检查用户是否存在
	sysUser, err := l.svcCtx.SysUserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, xerr.NewCodeError(xerr.ErrUserNotFound)
		}
		l.Errorf("查询用户失败，用户ID[%d]：%v", in.Id, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	if sysUser.DeletedAt.Valid {
		l.Errorf("用户[%d]已删除", in.Id)
		return nil, xerr.NewCodeError(xerr.ErrParamInvalid)
	}

	// 校验用户名
	if in.Username != "" {
		userByName, err := l.svcCtx.SysUserModel.FindOneByUsername(l.ctx, in.Username)
		if err != nil {
			l.Errorf("更新用户时校验用户名[%s]失败：%v", in.Username, err)
			return nil, xerr.NewCodeError(xerr.ErrInternal)
		}

		if userByName != nil && userByName.Id != in.Id {
			l.Errorf("用户[%s]已存在", in.Username)
			return nil, xerr.NewCodeError(xerr.ErrParamInvalid)
		}
	}
	// 更新字典类型信息
	updates := make(map[string]interface{})
	if in.HasUsername {
		updates["username"] = in.Username
	}
	if in.HasNickname {
		updates["nickname"] = in.Nickname
	}
	if in.HasEmail {
		updates["email"] = in.Email
	}
	if in.HasPhone {
		updates["phone"] = in.Phone
	}
	if in.HasStatus {
		updates["status"] = in.Status
	}
	if in.HasAvatar {
		updates["avatar"] = in.Avatar
	}
	if in.HasRemark {
		updates["remark"] = in.Remark
	}

	// 是否需要触达角色绑定逻辑
	needUpdateRoles := in.RoleIds != nil

	if len(updates) == 0 && !needUpdateRoles {
		l.Error("更新字段为空")
		return nil, xerr.NewCodeError(xerr.ErrParamInvalid)
	}

	// 角色合法性校验
	if needUpdateRoles && len(in.RoleIds) > 0 {
		roles, err := l.svcCtx.SysRoleModel.FindByIds(l.ctx, in.RoleIds)
		if err != nil {
			l.Errorf("角色合法性校验失败：%v", err)
			return nil, xerr.NewCodeError(xerr.ErrInternal)
		}
		if len(roles) != len(in.RoleIds) { // 严格匹配，避免传无效ID被忽略
			return nil, xerr.NewCodeError(xerr.ErrRoleNotFound)
		}
	}

	// 开启事务
	err = l.svcCtx.Orm.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		if err := l.svcCtx.SysUserModel.UpdateTrans(l.ctx, tx, in.Id, updates); err != nil {
			l.Errorf("更新用户[%d]失败：%v", in.Id, err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		if needUpdateRoles {
			if err := l.svcCtx.SysUserRoleModel.AssignRolesTrans(l.ctx, tx, in.Id, in.RoleIds); err != nil {
				l.Errorf("更新用户[%d]分配角色失败：%v", in.Id, err)
				return xerr.NewCodeError(xerr.ErrInternal)
			}

		}
		return err
	})

	if err != nil {
		l.Errorf("更新用户事务执行失败: %v", err)
		return &pb.Empty{}, err
	}

	return &pb.Empty{}, nil
}
