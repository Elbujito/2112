package proc

import (
	"fmt"

	"github.com/Elbujito/2112/src/app-service/internal/clients/dbc"
	"github.com/Elbujito/2112/src/app-service/internal/clients/logger"
)

func DBDrop() {
	logger.SetLogger(string(logger.DebugLvl))

	dbClient := dbc.GetDBClient()

	dbClient.InitServerConnection()

	if err := dbClient.DropDatabase(); err != nil {
		panic(fmt.Errorf("failed to drop database: %w", err))
	}

	// logger.Info("Database '" + config.Env.Config.DBName + "' dropped.")

}
