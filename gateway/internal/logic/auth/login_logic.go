// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"go-zero-rpc/sys-rpc/pb"
	"net/http"

	"go-zero-rpc/gateway/internal/svc"
	"go-zero-rpc/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq, r *http.Request) (resp *types.LoginResp, err error) {

	// 把 HTTP 上下文信息注入到 RPC 请求里
	rpcResp, err := l.svcCtx.AuthRpc.Login(r.Context(), &sys.LoginReq{
		Username:  req.Username,
		Password:  req.Password,
		ClientIp:  getClientIP(r),
		UserAgent: r.UserAgent(),
	})
	if err != nil {
		return nil, err
	}

	return &types.LoginResp{
		AccessToken:  rpcResp.AccessToken,
		RefreshToken: rpcResp.RefreshToken,
		ExpiresIn:    rpcResp.ExpiresIn,
		UserInfo: types.UserInfo{
			UserId:   rpcResp.UserInfo.UserId,
			Username: rpcResp.UserInfo.Username,
			Nickname: rpcResp.UserInfo.Nickname,
			Email:    rpcResp.UserInfo.Email,
			Phone:    rpcResp.UserInfo.Phone,
			Avatar:   rpcResp.UserInfo.Avatar,
			Roles:    rpcResp.UserInfo.Roles,
			Menus:    convertMenuItems(rpcResp.UserInfo.Menus),
		},
	}, err
}

func getClientIP(r *http.Request) string {
	// 复制 internal/middleware/auth_middleware.go 的 GetClientIP
	if v := r.Header.Get("X-Forwarded-For"); v != "" {
		return v
	}
	if v := r.Header.Get("X-Real-IP"); v != "" {
		return v
	}
	return r.RemoteAddr
}

func convertMenuItems(in []*sys.MenuItem) []types.MenuItem {
	out := make([]types.MenuItem, 0, len(in))
	for _, item := range in {
		if item == nil {
			continue
		}
		out = append(out, types.MenuItem{
			Id:        item.Id,
			ParentId:  item.ParentId,
			MenuName:  item.MenuName,
			MenuType:  item.MenuType,
			Path:      item.Path,
			Component: item.Component,
			Icon:      item.Icon,
			Sort:      item.Sort,
			Perms:     item.Perms,
			Status:    int(item.Status),
			Children:  convertMenuItems(item.Children),
		})
	}
	return out
}
