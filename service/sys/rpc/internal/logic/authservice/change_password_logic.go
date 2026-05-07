package authservicelogic

import (
	"context"
	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/pkg/encrypt"

	"go-zero-rpc/sys-rpc/internal/svc"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdatePassword 修改当前用户的登录密码。

/*
	业务流程：

//  1. 获取当前用户信息
//  2. 验证旧密码是否正确
//  3. 对新密码进行bcrypt加密
//  4. 更新数据库中的密码字段
//
// 参数：
//   - req : 修改密码请求体（旧密码 + 新密码）
//
// 返回：
//   - error : 业务错误
*/

func (l *ChangePasswordLogic) ChangePassword(in *sys.ChangePasswordReq) (*sys.CommonResp, error) {

	if in.UserId == 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}

	// 1. 查询用户信息（需要获取当前密码Hash）
	user, err := l.svcCtx.SysUserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, xerr.NewCodeError(xerr.ErrUserNotFound)
		}
		l.Logger.Errorf("查询用户[%d]失败：%v", in.UserId, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	if user.DeletedAt.Valid {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "用户已被删除")
	}
	// 2. 验证旧密码
	if !encrypt.CheckPassword(in.OldPassword, user.Password) {
		return nil, xerr.NewCodeError(xerr.ErrOldPasswordWrong)
	}

	// 3. 对新密码进行bcrypt加密
	newHashedPassword, err := encrypt.HashPassword(in.NewPassword)
	if err != nil {
		l.Logger.Errorf("新密码加密失败：%v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	// 4. 更新密码
	if err = l.svcCtx.SysUserModel.UpdatePassword(l.ctx, in.UserId, newHashedPassword); err != nil {
		l.Logger.Errorf("更新用户[%d]密码失败：%v", in.UserId, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	return &sys.CommonResp{
		Code: 0,
		Msg:  "密码成功更新",
	}, nil
}
