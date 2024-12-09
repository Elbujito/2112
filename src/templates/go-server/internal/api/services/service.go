package serviceapi

import (
	"github.com/Elbujito/2112/src/template/go-server/internal/config"
	"github.com/Elbujito/2112/src/template/go-server/internal/data"
	repository "github.com/Elbujito/2112/src/template/go-server/internal/repositories"
	"github.com/Elbujito/2112/src/template/go-server/internal/services"
)

// ServiceComponent holds all service instances for dependency injection.
type ServiceComponent struct {
	TestService services.TestService
}

// NewServiceComponent initializes and returns a new ServiceComponent.
func NewServiceComponent(env *config.SEnv) *ServiceComponent {
	// Initialize database connection
	database := data.NewDatabase()

	// Initialize repositories
	testRepo := repository.NewTestRepository(&database)

	// Create services
	testService := services.NewTestService(testRepo)

	return &ServiceComponent{
		TestService: testService,
	}
}
