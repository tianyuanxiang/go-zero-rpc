// Package encrypt 提供密码加密和验证工具函数。
// 使用bcrypt算法进行密码哈希，适合存储用户密码。
package encrypt

import (
	"golang.org/x/crypto/bcrypt"
)

// bcryptCost bcrypt加密强度，值越大越安全但越慢。
// 10是推荐的生产环境值，约需100ms计算时间。
const bcryptCost = 10

// HashPassword 对明文密码进行bcrypt加密，返回哈希字符串。
//
// 参数：
//   - password : 用户输入的明文密码
//
// 返回：
//   - string : bcrypt哈希字符串（含盐值，可直接存入数据库）
//   - error  : 加密失败时的错误信息
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword 验证明文密码与bcrypt哈希是否匹配。
//
// 参数：
//   - password : 用户输入的明文密码
//   - hash     : 数据库中存储的bcrypt哈希字符串
//
// 返回：
//   - bool : true表示密码正确，false表示密码错误
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
