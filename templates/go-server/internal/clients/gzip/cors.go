package gzip

import (
	"reflect"

	"github.com/Elbujito/2112/template/go-server/internal/config/features"
)

type GzipClient struct {
	name   string
	config features.GzipConfig
}

func (c *GzipClient) Name() string {
	return c.name
}

func (c *GzipClient) Configure(v any) {
	c.config = v.(reflect.Value).Interface().(features.GzipConfig)
}

func (c *GzipClient) GetConfig() features.GzipConfig {
	return c.config
}
