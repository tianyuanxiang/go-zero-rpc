package file

import (
	"context"

	"go-zero-rpc/common/middleware"
	"go-zero-rpc/common/xerr"
	"go-zero-rpc/gateway/internal/svc"
	sys "go-zero-rpc/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileLogic {
	return &DeleteFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteFileLogic) DeleteFile(fileID int64) error {
	userID := middleware.GetUserIdFromCtx(l.ctx)
	if userID == 0 {
		return xerr.NewCodeErrorMsg(xerr.ErrUnauthorized, "鏈巿鏉冿紝璇峰厛鐧诲綍")
	}

	_, err := l.svcCtx.SysRpc.DeleteFile(l.ctx, &sys.DeleteFileReq{
		FileId:     fileID,
		OperatorId: userID,
	})
	if err != nil {
		l.Errorf("delete file rpc failed, userId=%d, fileId=%d, err=%v", userID, fileID, err)
		return err
	}

	return nil
}
