package serviceapi

import (
	"github.com/Elbujito/2112/internal/clients/celestrack"
	propagator "github.com/Elbujito/2112/internal/clients/propagate"
	"github.com/Elbujito/2112/internal/data"
	repository "github.com/Elbujito/2112/internal/repositories"
	"github.com/Elbujito/2112/internal/services"
)

// ServiceComponent holds all service instances for dependency injection.
type ServiceComponent struct {
	SatelliteService services.SatelliteService
	TileService      services.TileService
}

// NewServiceComponent initializes and returns a new ServiceComponent.
func NewServiceComponent() *ServiceComponent {
	// Initialize database connection
	database := data.NewDatabase()

	// Initialize repositories
	tleRepo := repository.NewTLERepository(&database)
	satelliteRepo := repository.NewSatelliteRepository(&database)
	tileRepo := repository.NewTileRepository(&database)

	// Initialize external clients
	propagteClient := propagator.NewPropagatorClient(propagator.DefaultPropagationAPIURL)
	celestrackClient := celestrack.CelestrackClient{}

	// Create services
	satelliteService := services.NewSatelliteService(tleRepo, propagteClient, &celestrackClient, satelliteRepo)
	tileService := services.NewTileService(tileRepo)

	return &ServiceComponent{
		SatelliteService: satelliteService,
		TileService:      tileService,
	}
}
