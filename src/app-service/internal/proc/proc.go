package proc

import (
	"fmt"

	"github.com/Elbujito/2112/src/app-service/internal/api/routers"
	clientsPkg "github.com/Elbujito/2112/src/app-service/internal/clients"
	"github.com/Elbujito/2112/src/app-service/internal/clients/cors"
	"github.com/Elbujito/2112/src/app-service/internal/clients/dbc"
	"github.com/Elbujito/2112/src/app-service/internal/clients/gzip"
	"github.com/Elbujito/2112/src/app-service/internal/clients/kratos"
	"github.com/Elbujito/2112/src/app-service/internal/clients/logger"
	"github.com/Elbujito/2112/src/app-service/internal/clients/service"
	"github.com/Elbujito/2112/src/app-service/internal/config"
	"github.com/Elbujito/2112/src/app-service/internal/data/models"
	"github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xutils"
)

func InitServiceEnv(serviceName string, version string) {
	config.SetServiceName(serviceName)
	config.SetServiceVersion(version)
	config.InitEnv()
	config.ResolveDevMode()
	config.InitFeatures()
	config.ResolveFlags()
	config.PrintEnvInEnvMode()
}

var clients []clientsPkg.IClient

func InitClients() {
	InitServiceClient()
	InitCorsClient()
	InitGzipClient()
	InitDbClient()
	InitOryKratosClient()
	// ...
}

func ConfigureClients() {
	logger.Debug("Configuring clients ...")
	for _, c := range clients {
		feature := config.Feature(c.Name())
		if feature.IsEnabled() {
			logger.Debug("Configuring %s client ...", c.Name())
			c.Configure(feature.Config)
			continue
		}
		logger.Warn("Client: '%s' is disabled, This may cause runtime errors if this client is used.", c.Name())
	}
}

func addClient(client clientsPkg.IClient) {
	clients = append(clients, client)
}

func InitServiceClient() {
	client := service.GetClient()
	logger.Debug("Activating %s client ...", client.Name())
	addClient(client)
}

func InitCorsClient() {
	client := cors.GetClient()
	logger.Debug("Activating %s client ...", client.Name())
	addClient(client)
}

func InitGzipClient() {
	client := gzip.GetClient()
	logger.Debug("Activating %s client ...", client.Name())
	addClient(client)
}

func InitDbClient() {
	client := dbc.GetDBClient()
	logger.Debug("Activating %s client ...", client.Name())
	client.SetSilent(!config.DevModeFlag)
	addClient(client)
}

func InitOryKratosClient() {
	client := kratos.GetClient()
	logger.Debug("Activating %s client ...", client.Name())
	addClient(client)
}

func InitDbConnection() {
	logger.Debug("Initializing database connection ...")
	dbc.GetDBClient().InitDBConnection()
}

func InitModels() {
	logger.Debug("Activating models ...")
	models.Init(dbc.GetDBClient().DB)
}

func PrintHiddenRoutesTable() {
	routers.InitHiddenAPIRouter()
	routes := routers.HiddenAPIRouter().Echo.Routes()

	t := xutils.PrepareRoutesTable(routes, "Hidden API Routes")
	xutils.SetTableBorderStyle(t, config.NoBorderFlag)

	fmt.Println(t.Render())
}

func PrintProtectedRoutesTable() {
	routers.InitProtectedAPIRouter()
	routes := routers.ProtectedAPIRouter().Echo.Routes()

	t := xutils.PrepareRoutesTable(routes, "Protected API Routes")
	xutils.SetTableBorderStyle(t, config.NoBorderFlag)

	fmt.Println(t.Render())
}

func PrintPublicRoutesTable() {
	publicApiRouter := routers.InitPublicAPIRouter(config.Env)
	routes := publicApiRouter.Echo.Routes()

	t := xutils.PrepareRoutesTable(routes, "Public API Routes")
	xutils.SetTableBorderStyle(t, config.NoBorderFlag)

	fmt.Println(t.Render())
}
