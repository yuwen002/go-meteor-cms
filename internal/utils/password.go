package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// MinPasswordLength 密码最小长度
const MinPasswordLength = 6

// ValidatePassword 校验密码长度是否合法
func ValidatePassword(password string) bool {
	return len(password) >= MinPasswordLength
}

// HashPassword 对密码进行加密
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 校验明文密码与哈希是否匹配
func CheckPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
