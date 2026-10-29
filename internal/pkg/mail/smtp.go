package mail

import (
	"errors"
	"net/smtp"
)

// Config SMTP 配置
type Config struct {
	// Enabled 是否启用真实邮件发送，false 时降级为控制台输出
	Enabled bool
	// Host SMTP 服务器地址
	Host string
	// Port SMTP 端口，常用 465（SSL）/587（STARTTLS）/25
	Port int
	// Username 登录账号，通常就是发件邮箱
	Username string
	// Password 登录密码或授权码
	Password string
	// From 发件人地址
	From string
}

// SMTPSender 基于 net/smtp 的发送实现
type SMTPSender struct {
	Config Config
}

// NewSMTPSender 创建 SMTP 发送器
func NewSMTPSender(cfg Config) *SMTPSender {
	if cfg.Port == 0 {
		cfg.Port = 587
	}
	return &SMTPSender{Config: cfg}
}

// Send 发送邮件
func (s *SMTPSender) Send(msg Message) error {
	if len(msg.To) == 0 {
		return errors.New("收件人为空")
	}
	if s.Config.From == "" {
		return errors.New("未配置发件人地址")
	}

	auth := smtp.PlainAuth("", s.Config.Username, s.Config.Password, s.Config.Host)
	return smtp.SendMail(
		addr(s.Config.Host, s.Config.Port),
		auth,
		s.Config.From,
		msg.To,
		[]byte(buildMIME(s.Config.From, msg)),
	)
}
