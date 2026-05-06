// Package jwtx 提供JWT令牌的生成和验证工具函数。
// 支持访问令牌（AccessToken）和刷新令牌（RefreshToken）双Token机制。
package jwtx

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenType 令牌类型枚举
type TokenType string

const (
	// TokenTypeAccess 访问令牌类型
	TokenTypeAccess TokenType = "access"
	// TokenTypeRefresh 刷新令牌类型
	TokenTypeRefresh TokenType = "refresh"
)

// Claims JWT自定义声明结构体，继承标准声明。
type Claims struct {
	// UserId 用户ID
	UserId int64 `json:"userId"`
	// Username 用户名
	Username string `json:"username"`
	// TokenType 令牌类型（access/refresh）
	TokenType TokenType `json:"tokenType"`
	jwt.RegisteredClaims
}

// GenerateToken 生成访问令牌（AccessToken）。
//
// 参数：
//   - userId   : 用户ID
//   - username : 用户名
//   - secret   : JWT签名密钥
//   - expire   : 过期时间（秒）
//
// 返回：
//   - string : 生成的JWT字符串
//   - error  : 生成失败时的错误信息
func GenerateToken(userId int64, username string, secret string, expire int64) (string, error) {
	now := time.Now()
	claims := Claims{
		UserId:    userId,
		Username:  username,
		TokenType: TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expire) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// GenerateRefreshToken 生成刷新令牌（RefreshToken）。
//
// 刷新令牌用于在访问令牌过期后换取新的访问令牌，
// 过期时间通常比访问令牌长（如7天）。
//
// 参数：
//   - userId   : 用户ID
//   - username : 用户名
//   - secret   : JWT签名密钥
//   - expire   : 过期时间（秒）
//
// 返回：
//   - string : 生成的JWT字符串
//   - error  : 生成失败时的错误信息
func GenerateRefreshToken(userId int64, username string, secret string, expire int64) (string, error) {
	now := time.Now()
	claims := Claims{
		UserId:    userId,
		Username:  username,
		TokenType: TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expire) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ParseToken 解析并验证JWT令牌。
//
// 参数：
//   - tokenStr : JWT字符串
//   - secret   : JWT签名密钥
//
// 返回：
//   - *Claims : 解析出的Claims信息
//   - error   : 验证失败时的错误信息（包含过期、格式错误等）
func ParseToken(tokenStr string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法必须是HS256，防止算法切换攻击
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// GetTokenExpireAt 获取令牌的过期时间戳（Unix秒）。
//
// 参数：
//   - tokenStr : JWT字符串
//   - secret   : JWT签名密钥
//
// 返回：
//   - int64 : 过期时间Unix时间戳，解析失败时返回0
func GetTokenExpireAt(tokenStr string, secret string) int64 {
	claims, err := ParseToken(tokenStr, secret)
	if err != nil {
		return 0
	}
	if claims.ExpiresAt == nil {
		return 0
	}
	return claims.ExpiresAt.Unix()
}
