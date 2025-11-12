package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("无效或过期的 token")

// GenerateToken 生成 JWT Token
// 参数:
//   - secret: 用于签名的密钥
//   - expire: Token 过期时间（秒）
//   - payload: 需要存储在 Token 中的自定义数据
//
// 返回:
//   - string: 生成的 Token 字符串
//   - error: 生成过程中发生的错误
func GenerateToken(secret string, expire int64, payload map[string]interface{}) (string, error) {
	// 创建 JWT Claims 对象
	claims := jwt.MapClaims{}
	// 将自定义数据存入 Claims
	for k, v := range payload {
		claims[k] = v
	}
	// 设置 Token 过期时间
	claims["exp"] = time.Now().Add(time.Second * time.Duration(expire)).Unix()

	// 使用 HS256 算法创建 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用密钥对 Token 进行签名并返回
	return token.SignedString([]byte(secret))
}

// ParseToken 解析并验证 JWT Token
// 参数:
//   - secret: 用于验证签名的密钥，必须与生成 Token 时使用的密钥一致
//   - tokenStr: 待解析的 Token 字符串
//
// 返回:
//   - map[string]interface{}: 解析出的 Token 中的自定义数据
//   - error: 解析或验证失败时返回的错误
func ParseToken(secret, tokenStr string) (map[string]interface{}, error) {
	// 解析 Token 并验证签名
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// 返回用于验证签名的密钥
		return []byte(secret), nil
	})

	// 检查 Token 是否有效
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	// 将 claims 转换为 MapClaims 类型
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		data := make(map[string]interface{})
		// 复制除 exp 之外的所有声明到返回的 map 中
		for k, v := range claims {
			if k != "exp" { // 排除过期时间字段
				data[k] = v
			}
		}
		return data, nil
	}

	// 如果 claims 类型转换失败，返回无效 Token 错误
	return nil, ErrInvalidToken
}
