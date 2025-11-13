// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/yuwen002/go-meteor-cms/api/cms/internal/config"
	"github.com/yuwen002/go-meteor-cms/ent"
)

type ServiceContext struct {
	Config    config.Config
	EntClient *ent.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	client, err := ent.Open(c., "root:123456@tcp(127.0.0.1:3306)/go_meteor_cms?parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config: c,
		EntClient: client,
	}
}
