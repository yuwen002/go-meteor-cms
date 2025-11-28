// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/config"
	"github.com/yuwen002/go-meteor-cms/api/cms/v1/internal/middleware"
	"github.com/yuwen002/go-meteor-cms/ent"
	_ "github.com/yuwen002/go-meteor-cms/ent/runtime"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config        config.Config
	EntClient     *ent.Client
	JwtMiddleware rest.Middleware
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
	return &ServiceContext{
		Config:        c,
		EntClient:     client,
		JwtMiddleware: middleware.JwtMiddleware(&c),
	}
}
