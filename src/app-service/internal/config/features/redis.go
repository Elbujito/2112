package features

import xconstants "github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xconstants"

type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     string `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
}

var redis = &Feature{
	Name:       xconstants.FEATURE_REDIS,
	Config:     &RedisConfig{},
	enabled:    true,
	configured: false,
	ready:      false,
	requirements: []string{
		"Host",
		"Port",
		"Password",
	},
}

func init() {
	Features.Add(redis)
}
