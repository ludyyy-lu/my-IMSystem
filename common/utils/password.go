package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用 bcrypt 对密码进行哈希处理
// 返回哈希后的密码字符串和可能的错误
func HashPassword(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	return string(bytes), err
}

// CheckPasswordHash 检查密码是否与哈希值匹配
// 返回布尔值表示是否匹配
func CheckPasswordHash(pwd string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err == nil
}
