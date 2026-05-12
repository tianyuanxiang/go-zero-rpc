package systemservicelogic

import (
	"context"
	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/pb"
	"go-zero-rpc/sys-rpc/pkg/encrypt"

	"go-zero-rpc/sys-rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ResetPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ResetPasswordLogic) ResetPassword(in *pb.ResetPasswordReq) (*pb.Empty, error) {
	// 获取当前用户userId
	operatorId := in.OperatorId
	if operatorId == 0 {
		return nil, xerr.NewCodeError(xerr.ErrUnauthorized)
	}
	// 目标用户ID合法性兜底校验（handler 已校验，此处防御性兜底）
	if in.UserId <= 0 {
		return nil, xerr.NewCodeError(xerr.ErrParamInvalid)
	}
	// 禁止通过本接口重置自己的密码，必须引导至"修改密码"功能（需校验旧密码）
	if in.UserId == operatorId {
		l.Infof("操作者[%d]尝试通过重置接口修改自身密码，已拦截", operatorId)
		return nil, xerr.NewCodeErrorMsg(xerr.ErrForbidden, "请通过修改密码功能操作")
	}

	// 查询目标用户是否存在且未被软删除
	targetUser, err := l.svcCtx.SysUserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, xerr.NewCodeError(xerr.ErrUserNotFound)
		}
		l.Errorf("操作者[%d]查询目标用户[%d]失败：%v", operatorId, in.UserId, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}
	if targetUser.DeletedAt.Valid {
		l.Infof("操作者[%d]尝试重置已删除用户[%d]的密码，已拦截", operatorId, in.UserId)
		return nil, xerr.NewCodeErrorMsg(xerr.ErrForbidden, "该用户已被删除")
	}

	// 新密码 bcrypt 加密
	hashedPassword, err := encrypt.HashPassword(in.NewPassword)
	if err != nil {
		l.Errorf("新密码加密失败：%v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	// 更新目标用户的密码（注意此处必须用 req.Id，不能用 operatorId）
	if err := l.svcCtx.SysUserModel.UpdatePassword(l.ctx, in.UserId, hashedPassword); err != nil {
		l.Errorf("操作者[%d]重置用户[%d]密码失败：%v", operatorId, in.UserId, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	// 高敏操作审计日志（密码重置必须留痕，便于事后追溯）
	l.Infof("操作者[%d]成功重置用户[%d]的密码", operatorId, in.UserId)

	return &pb.Empty{}, nil
}
