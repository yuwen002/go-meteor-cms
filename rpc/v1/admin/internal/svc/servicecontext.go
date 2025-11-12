package svc

import (
	"context"
	"go-meteor-cms/ent"
	"go-meteor-cms/rpc/v1/admin/internal/config"
	"log"
)

type ServiceContext struct {
	Config    config.Config
	EntClient *ent.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	client, err := ent.Open(c.Database.Driver, c.Database.Source)
	if err != nil {
		log.Fatalf("failed opening connection to database: %v", err)
	}

	// ⚠️ 仅开发环境自动建表
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	return &ServiceContext{
		Config:    c,
		EntClient: client,
	}
}
