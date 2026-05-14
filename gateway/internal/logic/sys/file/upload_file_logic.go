package file

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"go-zero-rpc/common/middleware"
	"go-zero-rpc/common/xerr"
	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	"go-zero-rpc/gateway/pkg/upload"
	sys "go-zero-rpc/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadFileLogic {
	return &UploadFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadFileLogic) UploadFile(r *http.Request) (*types.FileInfoResp, error) {
	if err := r.ParseMultipartForm(l.svcCtx.Config.Upload.MaxSize); err != nil {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "涓婁紶鏂囦欢澶у皬瓒呭嚭闄愬埗")
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, "璇烽€夋嫨涓婁紶鏂囦欢")
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowExt(l.svcCtx.Config.Upload.AllowedExts, ext) {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrParamInvalid, fmt.Sprintf("unsupported file type: %s", ext))
	}

	userID := middleware.GetUserIdFromCtx(l.ctx)
	if userID == 0 {
		return nil, xerr.NewCodeErrorMsg(xerr.ErrUnauthorized, "闁哄牜浜濆鍧楀级閸愯法绀夐悹鍥у槻閸樻盯鎯傜拠鑼Э")
	}

	savedFile, err := upload.SaveFile(file, header, l.svcCtx.Config.Upload.Path)
	if err != nil {
		l.Errorf("save upload file failed, name=%s, err=%v", header.Filename, err)
		return nil, xerr.NewCodeErrorMsg(xerr.ErrInternal, "鏂囦欢淇濆瓨澶辫触")
	}

	rpcResp, err := l.svcCtx.SysRpc.RegisterFile(l.ctx, &sys.RegisterFileReq{
		FileName:   savedFile.OriginName,
		FileUrl:    savedFile.FileURL,
		FileSize:   savedFile.FileSize,
		MimeType:   savedFile.MimeType,
		UploaderId: userID,
	})
	if err != nil {
		_ = upload.RemoveFile(savedFile.FilePath)
		l.Errorf("register file rpc failed, userId=%d, url=%s, err=%v", userID, savedFile.FileURL, err)
		return nil, err
	}
	if rpcResp.File == nil {
		_ = upload.RemoveFile(savedFile.FilePath)
		return nil, xerr.NewCodeErrorMsg(xerr.ErrInternal, "鏂囦欢娉ㄥ唽澶辫触")
	}

	return &types.FileInfoResp{
		FileId:   rpcResp.File.FileId,
		FileName: rpcResp.File.FileName,
		FileURL:  rpcResp.File.FileUrl,
		FileSize: rpcResp.File.FileSize,
		MimeType: rpcResp.File.MimeType,
	}, nil
}

func allowExt(allowedExts []string, ext string) bool {
	if len(allowedExts) == 0 {
		return true
	}

	for _, allowed := range allowedExts {
		if strings.EqualFold(allowed, ext) {
			return true
		}
	}

	return false
}
