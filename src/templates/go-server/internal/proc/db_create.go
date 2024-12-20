package proc

import (
	"fmt"

	"github.com/Elbujito/2112/src/templates/go-server/internal/clients/dbc"
	"github.com/Elbujito/2112/src/templates/go-server/internal/clients/logger"
)

func DBCreate() {
	// init feature [database]
	logger.SetLogger(string(logger.DebugLvl))

	dbClient := dbc.GetDBClient()

	dbClient.InitServerConnection()

	if err := dbClient.CreateDatabase(); err != nil {
		panic(fmt.Errorf("failed to create database: %w", err))
	}

	// logger.Info("Database '" + config.Env.Config.DBName + "' created successfully.")

}
