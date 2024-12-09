package features

import xconstants "github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xconstants"

type CorsConfig struct {
	AllowOrigins  string `mapstructure:"CORS_ALLOW_ORIGINS"`
	AllowMethods  string `mapstructure:"CORS_ALLOW_METHODS"`
	AllowHeaders  string `mapstructure:"CORS_ALLOW_HEADERS"`
	MaxAge        int    `mapstructure:"CORS_MAX_AGE"`
	AllowCreds    bool   `mapstructure:"CORS_ALLOW_CREDENTIALS"`
	ExposeHeaders string `mapstructure:"CORS_EXPOSE_HEADERS"`
}

var cors = &Feature{
	Name:       xconstants.FEATURE_CORS,
	Config:     &CorsConfig{},
	enabled:    true,
	configured: false,
	ready:      false,
	requirements: []string{
		"AllowOrigins",
	},
}

func init() {
	Features.Add(cors)
}
