// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package middleware

import (
	"context"
	"fmt"
	"go-zero-rpc/common/jwtx"
	"go-zero-rpc/common/response"
	"go-zero-rpc/gateway/internal/config"
	"net/http"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

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

// AuthMiddleware JWT认证中间件。
//
// 从请求头 Authorization 中提取Bearer Token，验证其有效性，
// 并将解析出的 userId 和 username 写入请求上下文，供后续 Handler/Logic 使用。
//
// 验证失败（无token、格式错误、已过期）时直接返回 401 响应，不继续处理。
//
// 参数：
//   - cfg : 应用程序配置（需要JWT密钥）
//
// 返回：
//   - func(http.Handler) http.Handler : 标准中间件函数
func AuthMiddleware(cfg config.Config, rdb *redis.Client) func(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 从请求头获取Authorization字段
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.FailUnauthorized(w, r)
				return
			}

			// 检查并去除Bearer前缀
			// 标准格式：Authorization: Bearer <token>
			if !strings.HasPrefix(authHeader, "Bearer ") {
				response.Fail(w, r, response.CodeUnauthorized, "Authorization格式错误，应为：Bearer <token>")
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenStr == "" {
				response.FailUnauthorized(w, r)
				return
			}

			// 解析并验证Token
			claims, err := jwtx.ParseToken(tokenStr, cfg.Auth.AccessSecret)
			if err != nil {
				logx.WithContext(r.Context()).Infof("Token验证失败：%v", err)
				response.Fail(w, r, response.CodeTokenExpired, response.MsgTokenExpired)
				return
			}

			// 仅允许访问令牌（Access Token）通过此中间件
			// 刷新令牌只能用于 /api/auth/refresh 接口
			if claims.TokenType != jwtx.TokenTypeAccess {
				response.Fail(w, r, response.CodeUnauthorized, "请使用访问令牌访问此接口")
				return
			}

			// 新增：检查 Token 是否在黑名单（已登出）
			if rdb != nil {
				blacklistKey := fmt.Sprintf("%s%s", RedisTokenBlacklistPrefix, tokenStr)
				exists, redisErr := rdb.Exists(r.Context(), blacklistKey).Result()
				if redisErr == nil && exists > 0 {
					response.Fail(w, r, response.CodeUnauthorized, "Token已失效，请重新登录")
					return
				}
			}
			// 将用户信息写入Context，供后续Handler使用
			ctx := context.WithValue(r.Context(), ContextKeyUserId, claims.UserId)
			ctx = context.WithValue(ctx, ContextKeyUsername, claims.Username)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIdFromCtx 从Context中获取当前登录用户的ID。
//
// 参数：
//   - ctx : 请求上下文
//
// 返回：
//   - int64 : 用户ID，如果Context中没有则返回0
func GetUserIdFromCtx(ctx context.Context) int64 {
	userId, ok := ctx.Value(ContextKeyUserId).(int64)
	if !ok {
		return 0
	}
	return userId
}

// GetUsernameFromCtx 从Context中获取当前登录用户的用户名。
//
// 参数：
//   - ctx : 请求上下文
//
// 返回：
//   - string : 用户名，如果Context中没有则返回空字符串
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
//
// 参数：
//   - r : HTTP请求
//
// 返回：
//   - string : 客户端IP地址
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
