package systemservicelogic

import (
	"context"
	"path"
	"path/filepath"
	"strings"

	"go-zero-rpc/common/xerr"
	systemmodel "go-zero-rpc/sys-rpc/internal/model"
	"go-zero-rpc/sys-rpc/internal/svc"
	"go-zero-rpc/sys-rpc/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterFileLogic {
	return &RegisterFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterFileLogic) RegisterFile(in *sys.RegisterFileReq) (*sys.RegisterFileResp, error) {
	originName := strings.TrimSpace(in.FileName)
	fileURL := strings.TrimSpace(in.FileUrl)
	if originName == "" || fileURL == "" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "文件名和文件地址不能为空")
	}

	storedName := path.Base(fileURL)
	if storedName == "" || storedName == "." || storedName == "/" {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "文件地址无效")
	}

	mimeType := strings.TrimSpace(in.MimeType)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	filePath := filepath.FromSlash(strings.TrimPrefix(fileURL, "/"))
	fileId, err := l.svcCtx.SysFileModel.InsertFileReturningId(l.ctx, &systemmodel.SysFile{
		Filename:   storedName,
		OriginName: originName,
		FilePath:   filePath,
		FileUrl:    fileURL,
		FileSize:   in.FileSize,
		FileType:   mimeType,
		UploaderId: in.UploaderId,
	})
	if err != nil {
		l.Errorf("register file failed, url=%s, err=%v", fileURL, err)
		return nil, xerr.NewCodeError(xerr.ErrInternal)
	}

	return &sys.RegisterFileResp{
		File: &sys.FileInfo{
			FileId:   fileId,
			FileName: originName,
			FileUrl:  fileURL,
			FileSize: in.FileSize,
			MimeType: mimeType,
		},
	}, nil
}
