package dbc

import (
	"github.com/Elbujito/2112/src/app-service/internal/clients/dbc/adapters"
	xconstants "github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xconstants"

	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
)

var dbClient *DBClient

func init() {
	dbClient = &DBClient{
		name:    xconstants.FEATURE_DATABASE,
		adapter: adapters.Adapters,
		silent:  true,
		gormConfig: &gorm.Config{
			Logger: gLogger.Default.LogMode(gLogger.Silent),
		},
	}
}

// GetDBClient definition
func GetDBClient() *DBClient {
	return dbClient
}
