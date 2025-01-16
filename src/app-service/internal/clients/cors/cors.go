package cors

import (
	"reflect"

	"github.com/Elbujito/2112/src/app-service/internal/config/features"
)

// CorsClient definition
type CorsClient struct {
	name   string
	config features.CorsConfig
}

// Name getters
func (c *CorsClient) Name() string {
	return c.name
}

// Configure sets config
func (c *CorsClient) Configure(v any) {
	c.config = v.(reflect.Value).Interface().(features.CorsConfig)
}

// GetConfig getters
func (c *CorsClient) GetConfig() features.CorsConfig {
	return c.config
}
