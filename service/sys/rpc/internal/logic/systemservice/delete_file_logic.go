package systemservicelogic

import (
	"context"
	"errors"
	"go-zero-rpc/sys-rpc/pb"

	"os"

	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DeleteFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileLogic {
	return &DeleteFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteFileLogic) DeleteFile(in *pb.DeleteFileReq) (*pb.Empty, error) {
	file, err := l.svcCtx.SysFileModel.FindActiveById(l.ctx, in.FileId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "文件不存在")
		}
		l.Errorf("find file failed, fileId=%d, err=%v", in.FileId, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	if file.FilePath != "" {
		if err := os.Remove(file.FilePath); err != nil && !errors.Is(err, os.ErrNotExist) {
			l.Errorf("remove file failed, fileId=%d, path=%s, err=%v", in.FileId, file.FilePath, err)
			return nil, xerr.NewCodeError(xerr.ErrInternal)
		}
	}

	if err := l.svcCtx.SysFileModel.SoftDeleteFile(l.ctx, in.FileId); err != nil {
		if err == sqlx.ErrNotFound {
			return nil, xerr.NewCodeErrorMsg(xerr.ErrNotFound, "文件不存在")
		}
		l.Errorf("soft delete file failed, fileId=%d, err=%v", in.FileId, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	return &pb.Empty{}, nil
}
