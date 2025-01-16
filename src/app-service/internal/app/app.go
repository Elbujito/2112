package app

import (
	"context"

	"github.com/Elbujito/2112/src/app-service/internal/clients/logger"
	"github.com/Elbujito/2112/src/app-service/internal/config"
	"github.com/Elbujito/2112/src/app-service/internal/proc"
	"github.com/Elbujito/2112/src/app-service/internal/services"
)

// App struct encapsulates shared dependencies
type App struct {
	Services *services.ServiceComponent
	Version  string
}

func NewApp(ctx context.Context, serviceName string, version string) (App, error) {

	proc.InitServiceEnv(serviceName, version)
	proc.InitClients()
	proc.ConfigureClients()
	proc.InitDbConnection()
	proc.InitModels()

	logger.Debug("App instance initialized with services.")

	return App{
		Services: services.NewServiceComponent(config.Env),
		Version:  version,
	}, nil
}
