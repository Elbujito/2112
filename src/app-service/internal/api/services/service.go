package serviceapi

import (
	"github.com/Elbujito/2112/src/app-service/internal/clients/celestrack"
	propagator "github.com/Elbujito/2112/src/app-service/internal/clients/propagate"
	"github.com/Elbujito/2112/src/app-service/internal/config"
	"github.com/Elbujito/2112/src/app-service/internal/data"
	repository "github.com/Elbujito/2112/src/app-service/internal/repositories"
	"github.com/Elbujito/2112/src/app-service/internal/services"
)

// ServiceComponent holds all service instances for dependency injection.
type ServiceComponent struct {
	SatelliteService services.SatelliteService
	TileService      services.TileService
}

// NewServiceComponent initializes and returns a new ServiceComponent.
func NewServiceComponent(env *config.SEnv) *ServiceComponent {
	// Initialize database connection
	database := data.NewDatabase()

	// Initialize repositories
	tleRepo := repository.NewTLERepository(&database)
	satelliteRepo := repository.NewSatelliteRepository(&database)
	tileRepo := repository.NewTileRepository(&database)

	// Initialize external clients
	propagteClient := propagator.NewPropagatorClient(env)
	celestrackClient := celestrack.NewCelestrackClient(env)

	// Create services
	satelliteService := services.NewSatelliteService(tleRepo, propagteClient, celestrackClient, satelliteRepo)
	tileService := services.NewTileService(tileRepo)

	return &ServiceComponent{
		SatelliteService: satelliteService,
		TileService:      tileService,
	}
}
