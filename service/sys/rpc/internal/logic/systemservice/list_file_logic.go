package systemservicelogic

import (
	"context"

	"go-zero-rpc/common/xerr"
	"go-zero-rpc/sys-rpc/internal/svc"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFileLogic {
	return &ListFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListFileLogic) ListFile(in *sys.ListFileReq) (*sys.ListFileResp, error) {
	files, total, err := l.svcCtx.SysFileModel.List(l.ctx, int(in.Page), int(in.PageSize), in.Keyword)
	if err != nil {
		l.Errorf("list file failed, err=%v", err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	list := make([]*sys.FileInfo, 0, len(files))
	for _, file := range files {
		if file == nil {
			continue
		}

		list = append(list, &sys.FileInfo{
			FileId:   file.Id,
			FileName: file.OriginName,
			FileUrl:  file.FileUrl,
			FileSize: file.FileSize,
			MimeType: file.FileType,
		})
	}

	return &sys.ListFileResp{
		Total: total,
		List:  list,
	}, nil
}
