package kratos

import xconstants "github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xconstants"

func init() {
	kratosClient = &Kratos{
		name:    xconstants.FEATURE_ORY_KRATOS,
		Session: &KratosSession{},
	}
}

func GetClient() *Kratos {
	return kratosClient
}

// var Cli *Kratos
// var Session *KratosSession
// var Config *KratosConfig

// func init() {
// 	Config = &KratosConfig{}
// }

// func InitKratos(devMode bool) {
// 	Cli = &Kratos{}
// 	Session = &KratosSession{}
// 	publicConfig := oryKratos.NewConfiguration()
// 	if devMode {
// 		publicConfig.Debug = true
// 	}
// 	publicConfig.Servers = []oryKratos.ServerConfiguration{
// 		{
// 			URL: Config.KratosPublicURL,
// 		},
// 	}
// 	Cli.Public = oryKratos.NewAPIClient(publicConfig)
// 	adminConfig := oryKratos.NewConfiguration()
// 	if devMode {
// 		adminConfig.Debug = true
// 	}
// 	adminConfig.Servers = []oryKratos.ServerConfiguration{
// 		{
// 			URL: Config.KratosAdminURL,
// 		},
// 	}
// 	Cli.admin = oryKratos.NewAPIClient(adminConfig)
// }

// func GetKratosInstance() *Kratos {
// 	return Cli
// }
