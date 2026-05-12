package systemservicelogic

import (
	"context"
	"go-zero-rpc/common/xerr"
	systemmodel "go-zero-rpc/sys-rpc/internal/model"
	"go-zero-rpc/sys-rpc/pb"
	"go-zero-rpc/sys-rpc/pkg/encrypt"

	"go-zero-rpc/sys-rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

type CreateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateUser 创建新用户。
//
// 业务流程：
//  1. 检查用户名是否已存在
//  2. 对密码进行bcrypt加密
//  3. 插入用户记录
//  4. 分配初始角色（如果提供了roleIds）
//
// 参数：
//   - in : 创建用户请求体
//
// 返回：
//   - *pb.CreateUserResp : 携带新建用户ID，供调用方做后续审计、跳转详情等
//   - error : 业务错误，遵循 gRPC status 约定（成功为 nil）
func (l *CreateUserLogic) CreateUser(in *pb.CreateUserReq) (*pb.CreateUserResp, error) {
	// 1. 检查用户名唯一性
	existUser, err := l.svcCtx.SysUserModel.FindOneByUsername(l.ctx, in.Username)
	if err != nil && err != sqlx.ErrNotFound {
		l.Logger.Errorf("查询用户名[%s]是否存在失败：%v", in.Username, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}
	if existUser != nil {
		return nil, xerr.NewCodeError(xerr.ErrUsernameDuplicate)
	}
	// 2.对密码进行加密
	hashedPassword, err := encrypt.HashPassword(in.Password)
	if err != nil {
		l.Errorf("密码加密失败：%v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}
	// 默认状态为启用
	status := in.Status
	if status == 0 {
		status = 1
	}
	// 3.插入用户
	// 开启事务
	// 用闭包外变量捕获事务内生成的用户ID，事务成功后供 resp 返回
	var newUserId int64
	err = l.svcCtx.Orm.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		// 1.插入用户
		userId, err := l.svcCtx.SysUserModel.InsertUserTrans(l.ctx, tx, &systemmodel.SysUser{
			Username: in.Username,
			Password: hashedPassword,
			Nickname: in.Nickname,
			Email:    in.Email,
			Phone:    in.Phone,
			Avatar:   in.Avatar,
			Status:   status,
			Remark:   in.Remark,
		})
		if err != nil {
			l.Errorf("插入用户记录失败：%v", err)
			return xerr.NewCodeError(xerr.ErrInternal)
		}
		newUserId = userId
		// 2. 分配初始角色
		if len(in.RoleIds) > 0 {
			// 确认角色真实存在
			roles, err := l.svcCtx.SysRoleModel.FindByIds(l.ctx, in.RoleIds)
			if err != nil {
				l.Errorf("确认角色真实存在时查询角色信息失败：%v", err)
				return xerr.NewCodeError(xerr.ErrInternal)
			}
			if len(roles) <= 0 {
				l.Infof("关联的角色ID错误, %s", in.RoleIds)
				return xerr.NewCodeError(xerr.ErrRoleNotFound)
			}
			if err = l.svcCtx.SysUserRoleModel.AssignRolesTrans(l.ctx, tx, userId, in.RoleIds); err != nil {
				l.Errorf("为新用户[%d]分配角色失败：%v", userId, err)
				// 角色分配失败不影响用户创建成功，仅记录日志
			}
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateUserResp{UserId: newUserId}, nil
}
