package app

import (
	"context"

	"github.com/Elbujito/2112/src/app-service/internal/config"
	"github.com/Elbujito/2112/src/app-service/internal/proc"
	"github.com/Elbujito/2112/src/app-service/internal/services"
	logger "github.com/Elbujito/2112/src/app-service/pkg/log"
)

// App struct encapsulates shared dependencies
type App struct {
	Services *services.ServiceComponent
	Version  string
}

func NewApp(ctx context.Context, serviceName string, version string) (App, error) {
	logger.Infof("Initializing app: serviceName=%s, version=%s", serviceName, version)

	// Initialize service environment
	logger.Debug("Initializing service environment...")
	proc.InitServiceEnv(serviceName, version)
	logger.Info("Service environment initialized.")

	// Initialize clients
	logger.Debug("Initializing clients...")
	proc.InitClients()
	logger.Info("Clients initialized.")

	// Configure clients
	logger.Debug("Configuring clients...")
	proc.ConfigureClients()
	logger.Info("Clients configured.")

	// Initialize database connection
	logger.Debug("Initializing database connection...")
	proc.InitDbConnection()
	logger.Info("Database connection initialized.")

	// Initialize models
	logger.Debug("Initializing models...")
	proc.InitModels()
	logger.Info("Models initialized.")

	// Finalize app instance creation
	logger.Debug("Creating service component...")
	serviceComponent := services.NewServiceComponent(config.Env)
	logger.Info("Service component created.")

	logger.Infof("App instance successfully initialized for serviceName=%s, version=%s", serviceName, version)

	return App{
		Services: serviceComponent,
		Version:  version,
	}, nil
}
