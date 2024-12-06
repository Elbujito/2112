package features

import "github.com/Elbujito/2112/fx/constants"

type KratosConfig struct {
	PublicService string `mapstructure:"KRATOS_PUBLIC_SERVICE"`
	AdminService  string `mapstructure:"KRATOS_ADMIN_SERVICE"`
}

var kratos = &Feature{
	Name:       constants.FEATURE_ORY_KRATOS,
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
