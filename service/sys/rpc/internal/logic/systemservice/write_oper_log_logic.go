// Code scaffolded by goctl. Safe to edit.
package systemservicelogic

import (
	"context"
	"database/sql"
	"go-zero-rpc/sys-rpc/pb"
	"strconv"
	"time"

	"go-zero-rpc/common/xerr"
	systemmodel "go-zero-rpc/sys-rpc/internal/model"
	"go-zero-rpc/sys-rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type WriteOperLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWriteOperLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WriteOperLogLogic {
	return &WriteOperLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// WriteOperLog 写入一条操作日志。
//
// 由网关或其他服务在请求处理完成后调用，用于记录管理后台的关键操作。
//
// 字段说明：
//   - oper_type 在 proto 里以字符串形式承载业务类型枚举值（如 "1"=新增）；
//     由调用方编码，本服务原样写入 business_type 字段。
//   - req_param / resp_result 为可选字段，使用 sql.NullString 存入。
func (l *WriteOperLogLogic) WriteOperLog(in *pb.WriteOperLogReq) (*pb.WriteOperLogResp, error) {
	// 将 proto 中的 oper_type（string）转为 model 中的 BusinessType（int64）
	var businessType int64
	if in.OperType != "" {
		if v, err := strconv.ParseInt(in.OperType, 10, 64); err == nil {
			businessType = v
		}
	}

	operLog := &systemmodel.SysOperLog{
		Title:         in.Title,
		BusinessType:  businessType,
		Method:        in.Method,
		RequestMethod: in.ReqMethod,
		OperatorType:  1, // 默认后台用户
		OperatorName:  in.OperName,
		OperatorId:    in.OperUserId,
		OperUrl:       in.ReqUrl,
		OperIp:        in.Ip,
		OperParam: sql.NullString{
			String: in.ReqParam,
			Valid:  in.ReqParam != "",
		},
		JsonResult: sql.NullString{
			String: in.RespResult,
			Valid:  in.RespResult != "",
		},
		Status:   in.Status,
		ErrorMsg: in.ErrorMsg,
		OperTime: time.Now(),
	}

	logId, err := l.svcCtx.SysOperLogModel.InsertOperLogReturningId(l.ctx, operLog)
	if err != nil {
		l.Errorf("写入操作日志失败：%v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	return &pb.WriteOperLogResp{Id: logId}, nil
}
