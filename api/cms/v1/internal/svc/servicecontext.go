// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package svc

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/config"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/middleware"
	"github.com/yuwen002/go-meteor-cms/ent"
	_ "github.com/yuwen002/go-meteor-cms/ent/runtime"
	"github.com/yuwen002/go-meteor-cms/internal/pkg/mail"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config                 config.Config
	EntClient              *ent.Client
	JwtMiddleware          rest.Middleware
	PermissionMiddleware   rest.Middleware
	MailSender             mail.Sender
	OperationLogMiddleware rest.Middleware
	ClientInfoMiddleware   rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {

	client, err := ent.Open(
		c.Database.Driver,
		c.Database.Source,
		ent.Debug(),
	)
	if err != nil {
		panic(err)
	}

	// 未启用 SMTP 时降级为控制台输出，保证本地与未配置环境可用
	var sender mail.Sender = mail.ConsoleSender{}
	if c.Mail.Enabled {
		sender = mail.NewSMTPSender(mail.Config{
			Enabled:  true,
			Host:     c.Mail.Host,
			Port:     c.Mail.Port,
			Username: c.Mail.Username,
			Password: c.Mail.Password,
			From:     c.Mail.From,
		})
	}

	return &ServiceContext{
		Config:                 c,
		EntClient:              client,
		JwtMiddleware:          middleware.NewJwtMiddleware(&c, client),
		PermissionMiddleware:   middleware.NewPermissionMiddleware(&c, client),
		MailSender:             sender,
		OperationLogMiddleware: middleware.NewOperationLogMiddleware(c, client),
		ClientInfoMiddleware:   middleware.NewClientInfoMiddleware().Handle,
	}
}
