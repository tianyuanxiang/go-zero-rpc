// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"
	"go-zero-rpc/sys-rpc/pb"

	"time"

	"go-zero-rpc/common/xerr"
	systemmodel "go-zero-rpc/sys-rpc/internal/model"
	"go-zero-rpc/sys-rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLoginLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLoginLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLoginLogLogic {
	return &ListLoginLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListLoginLog 分页查询登录日志，支持关键词、状态、时间范围筛选。
//
// 参数说明：
//   - in.HasStatus = true 时按 in.Status 过滤；否则不按状态过滤（传 -1 表示不过滤）
func (l *ListLoginLogLogic) ListLoginLog(in *pb.ListLoginLogReq) (*pb.ListLoginLogResp, error) {
	// 默认状态为 -1（不过滤），仅在 HasStatus 时使用上层指定的 status
	status := -1
	if in.HasStatus {
		status = int(in.Status)
	}

	listReq := &systemmodel.LoginLogListReq{
		Page:     int(in.Page),
		PageSize: int(in.PageSize),
		Keyword:  in.Keyword,
		Status:   status,
	}

	// 解析时间范围（可选参数）
	const timeLayout = "2006-01-02 15:04:05"
	if in.StartTime != "" {
		t, err := time.ParseInLocation(timeLayout, in.StartTime, time.Local)
		if err == nil {
			listReq.StartTime = t
		}
	}
	if in.EndTime != "" {
		t, err := time.ParseInLocation(timeLayout, in.EndTime, time.Local)
		if err == nil {
			listReq.EndTime = t
		}
	}

	logs, total, err := l.svcCtx.SysLoginLogModel.List(l.ctx, listReq)
	if err != nil {
		l.Errorf("查询登录日志失败：%v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	list := make([]*pb.LoginLogItem, 0, len(logs))
	for _, log := range logs {
		list = append(list, &pb.LoginLogItem{
			Id:        log.Id,
			UserId:    log.UserId,
			Username:  log.Username,
			Ip:        log.Ip,
			Location:  log.Location,
			UserAgent: log.Browser,
			Os:        log.Os,
			Status:    log.Status,
			Msg:       log.Msg,
			LoginTime: log.LoginTime.Format(timeLayout),
		})
	}

	return &pb.ListLoginLogResp{
		Total: total,
		List:  list,
	}, nil
}
