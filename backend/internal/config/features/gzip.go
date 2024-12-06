package features

import "github.com/Elbujito/2112/fx/constants"

type GzipConfig struct {
	Level string `mapstructure:"GZIP_LEVEL"`
}

var gzip = &Feature{
	Name:       constants.FEATURE_GZIP,
	Config:     &GzipConfig{},
	enabled:    true,
	configured: false,
	ready:      false,
	requirements: []string{
		"Level",
	},
}

func init() {
	Features.Add(gzip)
}
