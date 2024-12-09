package proc

import (
	"github.com/Elbujito/2112/src/app-service/internal/clients/dbc"
	"github.com/Elbujito/2112/src/app-service/internal/clients/logger"
	"github.com/Elbujito/2112/src/app-service/internal/data/migrations"
)

func DBMigrate() {
	logger.SetLogger(string(logger.DebugLvl))

	dbClient := dbc.GetDBClient()

	dbClient.InitDBConnection()

	migrations.Init(dbClient.DB)

	if err := migrations.Migrate(); err != nil {
		logger.Error("Failed to apply migrations: %s", err)
	}

	logger.Info("Migrations applied successfully")

}
