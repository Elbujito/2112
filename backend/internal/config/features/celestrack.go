package features

import "github.com/Elbujito/2112/fx/constants"

type CelestrackConfig struct {
	BaseUrl string `mapstructure:"CELESTRACK_URL"`
	Satcat  string `mapstructure:"CELESTRACK_SATCAT_URL"`
}

var celestrack = &Feature{
	Name:       constants.FEATURE_CELESTRACK,
	Config:     &CelestrackConfig{},
	enabled:    true,
	configured: false,
	ready:      false,
	requirements: []string{
		"BaseUrl",
		"SatcatUrl",
	},
}

func init() {
	Features.Add(celestrack)
}
