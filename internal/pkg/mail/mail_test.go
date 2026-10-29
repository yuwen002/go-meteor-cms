package mail

import (
	"strings"
	"testing"
)

func TestConsoleSenderNeverFails(t *testing.T) {
	sender := ConsoleSender{}
	if err := sender.Send(Message{To: []string{"a@example.com"}, Subject: "hi", Body: "body"}); err != nil {
		t.Fatalf("控制台发送器不应失败: %v", err)
	}
}

// TestBuildMIME 邮件报文需包含 UTF-8 编码的主题与正文，避免中文乱码
func TestBuildMIME(t *testing.T) {
	msg := Message{
		To:      []string{"user@example.com"},
		Subject: "重置密码",
		Body:    "你好，请重置密码",
	}
	raw := buildMIME("noreply@example.com", msg)

	must := []string{
		"From: noreply@example.com",
		"To: user@example.com",
		"Subject: =?UTF-8?B?",
		"Content-Type: text/plain; charset=UTF-8",
		"Content-Transfer-Encoding: base64",
	}
	for _, want := range must {
		if !strings.Contains(raw, want) {
			t.Fatalf("邮件报文缺少 %q:\n%s", want, raw)
		}
	}

	// 中文不应以明文出现在报文里（主题用 RFC2047，正文用 base64）
	if strings.Contains(raw, "重置密码") || strings.Contains(raw, "你好") {
		t.Fatalf("中文必须进行编码后再放入报文:\n%s", raw)
	}

	// 头部与正文之间必须有空行
	if !strings.Contains(raw, "\r\n\r\n") {
		t.Fatal("邮件头部与正文之间缺少空行分隔")
	}
}

// TestSMTPSenderValidation SMTP 发送器需在参数缺失时明确报错，而不是静默失败
func TestSMTPSenderValidation(t *testing.T) {
	sender := NewSMTPSender(Config{Enabled: true, Host: "smtp.example.com", Port: 587})

	if err := sender.Send(Message{To: nil, Subject: "s", Body: "b"}); err == nil {
		t.Fatal("收件人为空应返回错误")
	}
	if err := sender.Send(Message{To: []string{"a@b.c"}, Subject: "s", Body: "b"}); err == nil {
		t.Fatal("未配置发件人地址应返回错误")
	}
}

func TestAddr(t *testing.T) {
	if got := addr("smtp.example.com", 587); got != "smtp.example.com:587" {
		t.Fatalf("addr = %q", got)
	}
}
