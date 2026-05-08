// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"
	"database/sql"
	"time"

	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/svc"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearLoginLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClearLoginLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearLoginLogLogic {
	return &ClearLoginLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ClearLoginLog 清空所有登录日志（软删除）。
//
// 注：proto 中 ClearLoginLogReq 为空，语义为「清空全部」。
// 通过批量将 deleted_at 设置为当前时间实现软删除。
func (l *ClearLoginLogLogic) ClearLoginLog(in *sys.ClearLoginLogReq) (*sys.Empty, error) {
	result := l.svcCtx.Orm.WithContext(l.ctx).Table("sys_login_log").
		Where("deleted_at IS NULL").
		Update("deleted_at", sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		})
	if result.Error != nil {
		l.Errorf("清空登录日志失败：%v", result.Error)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	return &sys.Empty{}, nil
}
