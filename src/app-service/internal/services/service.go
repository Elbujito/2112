package services

import (
	"log"
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/clients/celestrack"
	propagator "github.com/Elbujito/2112/src/app-service/internal/clients/propagate"
	"github.com/Elbujito/2112/src/app-service/internal/clients/redis"
	"github.com/Elbujito/2112/src/app-service/internal/config"
	"github.com/Elbujito/2112/src/app-service/internal/data"
	repository "github.com/Elbujito/2112/src/app-service/internal/repositories"
)

// ServiceComponent holds all service instances for dependency injection.
type ServiceComponent struct {
	SatelliteService SatelliteService
	TileService      TileService
	ContextService   ContextService
}

// NewServiceComponent initializes and returns a new ServiceComponent.
func NewServiceComponent(env *config.SEnv) *ServiceComponent {
	database := data.NewDatabase()

	redisClient, err := redis.NewRedisClient(config.Env)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	tleRepo := repository.NewTLERepository(&database, redisClient, 24*7*time.Hour)
	satelliteRepo := repository.NewSatelliteRepository(&database)
	tileRepo := repository.NewTileRepository(&database)
	mappingRepo := repository.NewTileSatelliteMappingRepository(&database)
	contextRepo := repository.NewContextRepository(&database)

	propagteClient := propagator.NewPropagatorClient(env)
	celestrackClient := celestrack.NewCelestrackClient(env)

	satelliteService := NewSatelliteService(tleRepo, propagteClient, celestrackClient, satelliteRepo)
	tileService := NewTileService(tileRepo, tleRepo, satelliteRepo, mappingRepo)
	contextService := NewContextService(contextRepo)

	return &ServiceComponent{
		SatelliteService: satelliteService,
		TileService:      tileService,
		ContextService:   contextService,
	}
}
