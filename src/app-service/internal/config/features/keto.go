package features

import xconstants "github.com/Elbujito/2112/lib/fx/xconstants"

type KetoConfig struct {
	ReadService  string `mapstructure:"KETO_READ_SERVICE"`
	WriteService string `mapstructure:"KETO_WRITE_SERVICE"`
}

var keto = &Feature{
	Name:       xconstants.FEATURE_ORY_KETO,
	Config:     &KetoConfig{},
	enabled:    true,
	configured: false,
	ready:      false,
	requirements: []string{
		"ReadService",
		"WriteService",
	},
}

func init() {
	Features.Add(keto)
}
