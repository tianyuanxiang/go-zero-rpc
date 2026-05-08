// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package file

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-rpc/gateway/internal/svc"
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

func (l *DeleteFileLogic) DeleteFile(fileId int64) error {
	// TODO: 文件删除RPC（RegisterFile/DeleteFile/ListFile）尚未实现，
	// proto 中对应的 rpc 定义已注释，待后续开启后补全此逻辑。

	return nil
}
