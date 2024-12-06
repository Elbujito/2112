package features

import xconstants "github.com/Elbujito/2112/lib/fx/xconstants"

type KratosConfig struct {
	PublicService string `mapstructure:"KRATOS_PUBLIC_SERVICE"`
	AdminService  string `mapstructure:"KRATOS_ADMIN_SERVICE"`
}

var kratos = &Feature{
	Name:       xconstants.FEATURE_ORY_KRATOS,
	Config:     &KratosConfig{},
	enabled:    true,
	configured: false,
	ready:      false,
	requirements: []string{
		"PublicService",
		"AdminService",
	},
}

func init() {
	Features.Add(kratos)
}
