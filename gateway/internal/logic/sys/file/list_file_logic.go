package file

import (
	"context"

	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"
	sys "go-zero-rpc/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFileLogic {
	return &ListFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListFileLogic) ListFile(req *types.ListFileReq) (*types.ListFileResp, error) {
	rpcResp, err := l.svcCtx.SysRpc.ListFile(l.ctx, &sys.ListFileReq{
		Page:     int64(req.Page),
		PageSize: int64(req.PageSize),
		Keyword:  req.Keyword,
	})
	if err != nil {
		l.Errorf("list file rpc failed, err=%v", err)
		return nil, err
	}

	list := make([]types.FileInfoResp, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		if item == nil {
			continue
		}

		list = append(list, types.FileInfoResp{
			FileId:   item.FileId,
			FileName: item.FileName,
			FileURL:  item.FileUrl,
			FileSize: item.FileSize,
			MimeType: item.MimeType,
		})
	}

	return &types.ListFileResp{
		Total: rpcResp.Total,
		List:  list,
	}, nil
}
