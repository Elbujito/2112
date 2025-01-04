package routers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/api/handlers/errors"
	healthHandlers "github.com/Elbujito/2112/src/app-service/internal/api/handlers/healthz"
	"github.com/Elbujito/2112/src/app-service/internal/api/handlers/satellites"
	"github.com/Elbujito/2112/src/app-service/internal/api/handlers/tiles"
	"github.com/Elbujito/2112/src/app-service/internal/api/middlewares"
	serviceapi "github.com/Elbujito/2112/src/app-service/internal/api/services"
	"github.com/Elbujito/2112/src/app-service/internal/clients/logger"
	"github.com/Elbujito/2112/src/app-service/internal/config"
	xconstants "github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xconstants"

	"github.com/labstack/echo/v4"
)

// PublicRouter manages the public API router and its dependencies.
type PublicRouter struct {
	Echo             *echo.Echo
	Name             string
	ServiceComponent *serviceapi.ServiceComponent // Add ServiceComponent to Router
}

// Init initializes the Echo instance for the router.
func (r *PublicRouter) Init() {
	r.Echo = echo.New()
	r.Echo.HideBanner = true
	r.Echo.Logger = logger.GetLogger()
}

// InitPublicAPIRouter initializes and returns the public API router.
func InitPublicAPIRouter(env *config.SEnv) *PublicRouter {
	logger.Debug("Initializing public API router ...")

	// Initialize ServiceComponent
	serviceComponent := serviceapi.NewServiceComponent(env)

	// Create and initialize PublicRouter
	publicApiRouter := &PublicRouter{
		Name:             "public API",
		ServiceComponent: serviceComponent,
	}
	publicApiRouter.Init()

	// Register middlewares, routes, and error handlers
	if config.DevModeFlag {
		publicApiRouter.registerPublicApiDevModeMiddleware()
	}
	publicApiRouter.registerPublicAPIMiddlewares()
	publicApiRouter.registerPublicApiHealthCheckHandlers()
	publicApiRouter.registerPublicApiSecurityMiddlewares()
	publicApiRouter.registerPublicAPIRoutes()
	publicApiRouter.registerPublicApiErrorHandlers()

	logger.Debug("Public API registration complete.")
	return publicApiRouter
}

// RegisterPreMiddleware registers a pre-middleware.
func (r *PublicRouter) RegisterPreMiddleware(middleware echo.MiddlewareFunc) {
	r.Echo.Pre(middleware)
}

// RegisterMiddleware registers a middleware.
func (r *PublicRouter) RegisterMiddleware(middleware echo.MiddlewareFunc) {
	r.Echo.Use(middleware)
}

// Start starts the Echo server with graceful shutdown.
func (r *PublicRouter) Start(host string, port string) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Start server
	go func() {
		r.Echo.Logger.Info(fmt.Sprintf("Starting %s server on port: %s", r.Name, port))
		if err := r.Echo.Start(host + ":" + port); err != nil && err != http.ErrServerClosed {
			r.Echo.Logger.Fatal(err)
			r.Echo.Logger.Fatal(xconstants.MSG_SERVER_SHUTTING_DOWN)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 20 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := r.Echo.Shutdown(ctx); err != nil {
		r.Echo.Logger.Fatal(err)
	}
}

// Register middlewares
func (r *PublicRouter) registerPublicAPIMiddlewares() {
	r.RegisterPreMiddleware(middlewares.SlashesMiddleware())
	r.RegisterMiddleware(middlewares.LoggerMiddleware())
	r.RegisterMiddleware(middlewares.TimeoutMiddleware())
	r.RegisterMiddleware(middlewares.RequestHeadersMiddleware())
	r.RegisterMiddleware(middlewares.ResponseHeadersMiddleware())

	if config.Feature(xconstants.FEATURE_GZIP).IsEnabled() {
		r.RegisterMiddleware(middlewares.GzipMiddleware())
	}
}

// Register development mode middlewares
func (r *PublicRouter) registerPublicApiDevModeMiddleware() {
	r.RegisterMiddleware(middlewares.BodyDumpMiddleware())
}

// Register security-related middlewares
func (r *PublicRouter) registerPublicApiSecurityMiddlewares() {
	r.RegisterMiddleware(middlewares.XSSCheckMiddleware())

	if config.Feature(xconstants.FEATURE_CORS).IsEnabled() {
		r.RegisterMiddleware(middlewares.CORSMiddleware())
	}
}

// Register health check handlers
func (r *PublicRouter) registerPublicApiHealthCheckHandlers() {
	health := r.Echo.Group("/health")
	health.GET("/alive", healthHandlers.Index)
	health.GET("/ready", healthHandlers.Ready)
}

// Register error handlers
func (r *PublicRouter) registerPublicApiErrorHandlers() {
	r.Echo.HTTPErrorHandler = errors.AutomatedHttpErrorHandler()
	r.Echo.RouteNotFound("/*", errors.NotFound)
}

// Register public API routes
func (r *PublicRouter) registerPublicAPIRoutes() {

	// Initialize the SatelliteHandler with the SatelliteService from ServiceComponent
	satelliteHandler := satellites.NewSatelliteHandler(r.ServiceComponent.SatelliteService)
	tileHandler := tiles.NewTileHandler(r.ServiceComponent.TileService)

	satellite := r.Echo.Group("/satellites")
	satellite.GET("/orbit", satelliteHandler.GetSatellitePositionsByNoradID)
	satellite.GET("/paginated", satelliteHandler.GetPaginatedSatellites)
	satellite.GET("/paginated/tles", satelliteHandler.GetPaginatedSatelliteInfo)

	// Initialize the SatelliteHandler with the SatelliteService from ServiceComponent
	tile := r.Echo.Group("/tiles")
	tile.GET("/all", tileHandler.GetAllTiles)
	tile.GET("/region", tileHandler.GetTilesInRegionHandler)
	tile.GET("/mappings", tileHandler.GetPaginatedSatelliteMappings)
	tile.PUT("/mappings/recompute/bynoradID", tileHandler.RecomputeMappingsByNoradID)
	tile.GET("/mappings/bynoradID", tileHandler.GetSatelliteMappingsByNoradID)
}
