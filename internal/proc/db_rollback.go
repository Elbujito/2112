package proc

import (
	"github.com/Elbujito/2112/internal/data/migrations"
	"github.com/Elbujito/2112/pkg/clients/dbc"
	"github.com/Elbujito/2112/pkg/clients/logger"
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
