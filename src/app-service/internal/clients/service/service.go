package service

import (
	"reflect"

	"github.com/Elbujito/2112/src/app-service/internal/config/features"
)

type ServiceClient struct {
	name   string
	config features.ServiceConfig
}

func (c *ServiceClient) Name() string {
	return c.name
}

func (c *ServiceClient) Configure(v any) {
	c.config = v.(reflect.Value).Interface().(features.ServiceConfig)
}

func (c *ServiceClient) GetConfig() features.ServiceConfig {
	return c.config
}
