package proc

import (
	"github.com/Elbujito/2112/internal/clients/dbc"
	"github.com/Elbujito/2112/internal/clients/logger"
	"github.com/Elbujito/2112/internal/data/migrations"
)

func DBRollback() {
	logger.SetLogger(string(logger.DebugLvl))

	dbClient := dbc.GetDBClient()

	dbClient.InitDBConnection()

	migrations.Init(dbClient.DB)

	if err := migrations.Rollback(); err != nil {
		logger.Error("Failed to rollback migrations: %s", err)
	}

}
