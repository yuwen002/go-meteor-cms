package utils

import (
	"github.com/mojocn/base64Captcha"
)

// 配置验证码
var captchaStore = base64Captcha.DefaultMemStore

// GenerateCaptcha 生成验证码
func GenerateCaptcha() (string, string, error) {
	// 配置数字验证码：高80、宽240、4位数字、干扰0.7、字体80大小
	driver := base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, captchaStore)
	id, b64s, _, err := cp.Generate()
	if err != nil {
		return "", "", err
	}
	return id, b64s, nil
}

// VerifyCaptcha 验证验证码
func VerifyCaptcha(id, answer string) bool {
	if id == "" || answer == "" {
		return false
	}
	return captchaStore.Verify(id, answer, true) // true 表示验证后清除
}
