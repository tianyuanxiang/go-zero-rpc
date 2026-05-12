package systemservicelogic

import (
	"context"
	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/pb"

	"go-zero-rpc/sys-rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *pb.GetUserReq) (*pb.UserItem, error) {
	// 查询用户基本信息
	userInfo, err := l.svcCtx.SysUserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, xerr.NewCodeError(xerr.ErrUserNotFound)
		}
		l.Errorf("查询用户[%d]失败：%v", in.UserId, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	if userInfo.DeletedAt.Valid {
		l.Errorf("用户已删除")
		return nil, xerr.NewCodeError(xerr.ErrUserNotFound)
	}
	// 查询用户角色
	roleIds, err := l.svcCtx.SysUserRoleModel.GetRoleIdsByUserId(l.ctx, in.UserId)
	if err != nil {
		l.Errorf("查询用户[%d]的角色失败：%v", in.UserId, err)
		roleIds = []int64{}
	}

	// 4. 批量查询所有角色信息
	roleCodes := make([]string, 0, len(roleIds)) // roleId -> code
	if len(roleIds) > 0 {
		roles, err := l.svcCtx.SysRoleModel.FindByIds(l.ctx, roleIds)
		if err != nil {
			l.Logger.Errorf("批量查询角色信息失败: %v", err)
			return nil, xerr.NewCodeError(xerr.ErrInternal)
		}
		for _, role := range roles {
			roleCodes = append(roleCodes, role.Code)
		}
	}

	return &pb.UserItem{
		Id:        userInfo.Id,
		Username:  userInfo.Username,
		Nickname:  userInfo.Nickname,
		Email:     userInfo.Email,
		Phone:     userInfo.Phone,
		Avatar:    userInfo.Avatar,
		Status:    userInfo.Status,
		Remark:    userInfo.Remark,
		Roles:     roleCodes,
		RoleIds:   roleIds,
		CreatedAt: userInfo.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
