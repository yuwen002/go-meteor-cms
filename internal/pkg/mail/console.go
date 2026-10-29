package mail

import (
	"fmt"
	"strings"
)

// ConsoleSender 不做真实投递，仅输出日志。
// 用于本地开发或未配置 SMTP 的部署环境。
type ConsoleSender struct{}

// Send 打印邮件内容到日志
func (ConsoleSender) Send(msg Message) error {
	fmt.Printf("📧 [mail/console] 收件人: %s\n主题: %s\n正文:\n%s\n",
		strings.Join(msg.To, ","),
		msg.Subject,
		msg.Body,
	)
	return nil
}
