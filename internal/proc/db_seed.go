package proc

import (
	"github.com/Elbujito/2112/internal/data/seeds"
	"github.com/Elbujito/2112/pkg/clients/dbc"
	"github.com/Elbujito/2112/pkg/clients/logger"
)

func DBSeed() {
	logger.SetLogger(string(logger.DebugLvl))

	dbClient := dbc.GetDBClient()

	dbClient.InitDBConnection()

	seeds.Init(dbClient.DB)

	if err := seeds.Apply(); err != nil {
		logger.Error("Failed to apply seeds: %s", err)
	}

	logger.Info("Seeds applied successfully")

}
