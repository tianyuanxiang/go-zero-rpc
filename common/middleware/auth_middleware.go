package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"go-zero-rpc/common/jwtx"
	"go-zero-rpc/common/response"
	permclient "go-zero-rpc/sys-rpc/client/permissionservice"
	"go-zero-rpc/sys-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

// AuthMiddleware JWT认证中间件。
type AuthMiddleware struct {
	permRpc permclient.PermissionService
}

// NewAuthMiddleware 创建认证中间件。
func NewAuthMiddleware(permRpc permclient.PermissionService) *AuthMiddleware {
	return &AuthMiddleware{permRpc: permRpc}
}

// ContextKey 自定义Context键类型，避免与其他包的键冲突。
type ContextKey string

const (
	// ContextKeyUserId Context中存储用户ID的键
	ContextKeyUserId ContextKey = "userId"
	// ContextKeyUsername Context中存储用户名的键
	ContextKeyUsername ContextKey = "username"
	// RedisTokenBlacklistPrefix Token黑名单在Redis中的键前缀
	RedisTokenBlacklistPrefix = "token:blacklist:"
)

// Handle 返回JWT认证中间件函数。
//
// 从请求头 Authorization 中提取Bearer Token，验证其有效性，
// 并将解析出的 userId 和 username 写入请求上下文。
//
// 参数：
//   - accessSecret: JWT访问令牌的密钥
func (m *AuthMiddleware) Handle(accessSecret string) func(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.FailUnauthorized(w, r)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				response.Fail(w, r, response.CodeUnauthorized, "Authorization格式错误，应为：Bearer <token>")
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenStr == "" {
				response.FailUnauthorized(w, r)
				return
			}

			claims, err := jwtx.ParseToken(tokenStr, accessSecret)
			if err != nil {
				logx.WithContext(r.Context()).Infof("Token验证失败：%v", err)
				response.Fail(w, r, response.CodeTokenExpired, response.MsgTokenExpired)
				return
			}

			if claims.TokenType != jwtx.TokenTypeAccess {
				response.Fail(w, r, response.CodeUnauthorized, "请使用访问令牌访问此接口")
				return
			}

			blacklistKey := fmt.Sprintf("%s%s", RedisTokenBlacklistPrefix, tokenStr)
			isBlack, err := m.permRpc.IsTokenRevoked(r.Context(), &pb.IsBlackListReq{BlacklistKey: blacklistKey})
			if err != nil {
				logx.WithContext(r.Context()).Errorf("Token list check fail., key=%s, err=%v", blacklistKey, err)
				response.Fail(w, r, response.CodeInternalError, "Token 校验失败，请稍后重试")
				return
			}

			if isBlack.IsBlack {
				logx.WithContext(r.Context()).Infof("Token is logout, refuse! key=%s", blacklistKey)
				response.Fail(w, r, response.CodeUnauthorized, "Token已失效，请重新登录")
				return
			}

			ctx := context.WithValue(r.Context(), ContextKeyUserId, claims.UserId)
			ctx = context.WithValue(ctx, ContextKeyUsername, claims.Username)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIdFromCtx 从Context中获取当前登录用户的ID。
func GetUserIdFromCtx(ctx context.Context) int64 {
	userId, ok := ctx.Value(ContextKeyUserId).(int64)
	if !ok {
		return 0
	}
	return userId
}

// GetUsernameFromCtx 从Context中获取当前登录用户的用户名。
func GetUsernameFromCtx(ctx context.Context) string {
	username, ok := ctx.Value(ContextKeyUsername).(string)
	if !ok {
		return ""
	}
	return username
}

// GetClientIP 从HTTP请求中获取客户端真实IP地址。
//
// 优先从 X-Forwarded-For 头获取，其次 X-Real-IP，最后从 RemoteAddr 截取。
func GetClientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		parts := splitString(forwarded, ",")
		if len(parts) > 0 {
			return trimSpace(parts[0])
		}
	}

	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	remoteAddr := r.RemoteAddr
	if idx := lastIndexByte(remoteAddr, ':'); idx > 0 {
		return remoteAddr[:idx]
	}

	return remoteAddr
}

// splitString 按分隔符分割字符串。
func splitString(s, sep string) []string {
	if sep == "" {
		return []string{s}
	}
	result := make([]string, 0)
	start := 0
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
		}
	}
	result = append(result, s[start:])
	return result
}

// trimSpace 去除字符串首尾空格。
func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

// lastIndexByte 查找字节在字符串中最后出现的位置。
func lastIndexByte(s string, c byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == c {
			return i
		}
	}
	return -1
}
