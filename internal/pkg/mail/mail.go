// Package mail 提供邮件发送能力。
//
// 设计要点：
//   - Sender 是统一抽象，便于测试替换与后续接入其他邮件服务；
//   - 未配置 SMTP（Enabled=false）时降级为 ConsoleSender，只打日志，保证本地开箱即用；
//   - 统一使用 UTF-8 + Base64 编码主题与正文，避免中文乱码。
package mail

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// Message 待发送的邮件
type Message struct {
	// To 收件人列表
	To []string
	// Subject 邮件主题
	Subject string
	// Body 纯文本正文
	Body string
}

// Sender 邮件发送抽象
type Sender interface {
	// Send 发送邮件
	Send(msg Message) error
}

// encodeRFC2047 按 RFC 2047 编码邮件主题
func encodeRFC2047(s string) string {
	return "=?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte(s)) + "?="
}

// buildMIME 构造完整邮件报文
func buildMIME(from string, msg Message) string {
	var b strings.Builder
	b.WriteString("From: " + from + "\r\n")
	b.WriteString("To: " + strings.Join(msg.To, ",") + "\r\n")
	b.WriteString("Subject: " + encodeRFC2047(msg.Subject) + "\r\n")
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	b.WriteString("Content-Transfer-Encoding: base64\r\n")
	b.WriteString("\r\n")

	encoded := base64.StdEncoding.EncodeToString([]byte(msg.Body))
	// 按 76 字符折行，符合 MIME 规范
	for i := 0; i < len(encoded); i += 76 {
		end := i + 76
		if end > len(encoded) {
			end = len(encoded)
		}
		b.WriteString(encoded[i:end] + "\r\n")
	}

	return b.String()
}

// addr 生成 SMTP 服务地址
func addr(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}
