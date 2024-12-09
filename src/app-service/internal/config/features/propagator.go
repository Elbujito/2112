package features

import xconstants "github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xconstants"

type PropagatorConfig struct {
	BaseUrl string `mapstructure:"PROPAGATOR_URL"`
}

var propagator = &Feature{
	Name:       xconstants.FEATURE_PROPAGATOR,
	Config:     &PropagatorConfig{},
	enabled:    true,
	configured: false,
	ready:      false,
	requirements: []string{
		"BaseUrl",
	},
}

func init() {
	Features.Add(propagator)
}
