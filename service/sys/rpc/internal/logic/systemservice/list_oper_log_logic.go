// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"
	"go-zero-rpc/sys-rpc/pb"

	"strconv"
	"time"

	"go-zero-rpc/common/xerr"
	systemmodel "go-zero-rpc/sys-rpc/internal/model"
	"go-zero-rpc/sys-rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOperLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListOperLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOperLogLogic {
	return &ListOperLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListOperLog 分页查询操作日志，支持关键词、业务类型、状态、时间范围筛选。
func (l *ListOperLogLogic) ListOperLog(in *pb.ListOperLogReq) (*pb.ListOperLogResp, error) {
	// 默认状态为 -1（不过滤），仅在 HasStatus 时使用上层指定的 status
	status := -1
	if in.HasStatus {
		status = int(in.Status)
	}

	listReq := &systemmodel.OperLogListReq{
		Page:         int(in.Page),
		PageSize:     int(in.PageSize),
		Keyword:      in.Keyword,
		BusinessType: in.BusinessType,
		Status:       status,
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

	logs, total, err := l.svcCtx.SysOperLogModel.List(l.ctx, listReq)
	if err != nil {
		l.Errorf("查询操作日志失败：%v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	list := make([]*pb.OperLogItem, 0, len(logs))
	for _, log := range logs {
		list = append(list, &pb.OperLogItem{
			Id:         log.Id,
			Title:      log.Title,
			OperType:   strconv.FormatInt(log.BusinessType, 10),
			Method:     log.Method,
			ReqMethod:  log.RequestMethod,
			OperName:   log.OperatorName,
			DeptName:   log.DeptName,
			ReqUrl:     log.OperUrl,
			ReqParam:   log.OperParam.String,
			RespResult: log.JsonResult.String,
			Status:     log.Status,
			Ip:         log.OperIp,
			OperTime:   log.OperTime.Format(timeLayout),
		})
	}

	return &pb.ListOperLogResp{
		Total: total,
		List:  list,
	}, nil
}
