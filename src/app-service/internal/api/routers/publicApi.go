package routers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	apicontext "github.com/Elbujito/2112/src/app-service/internal/api/handlers/context"
	"github.com/Elbujito/2112/src/app-service/internal/api/handlers/errors"
	healthHandlers "github.com/Elbujito/2112/src/app-service/internal/api/handlers/healthz"
	"github.com/Elbujito/2112/src/app-service/internal/api/handlers/satellites"
	"github.com/Elbujito/2112/src/app-service/internal/api/handlers/tiles"
	"github.com/Elbujito/2112/src/app-service/internal/api/middlewares"
	"github.com/Elbujito/2112/src/app-service/internal/clients/logger"
	"github.com/Elbujito/2112/src/app-service/internal/config"
	"github.com/Elbujito/2112/src/app-service/internal/services"
	"github.com/labstack/echo/v4"
)

// PublicRouter manages the public API router and its dependencies.
type PublicRouter struct {
	Echo             *echo.Echo
	Name             string
	ServiceComponent *services.ServiceComponent
}

// Init initializes the Echo instance for the router.
func (r *PublicRouter) Init() {
	r.Echo = echo.New()
	r.Echo.HideBanner = true
	r.Echo.Logger = logger.GetLogger()

}

// InitPublicAPIRouter initializes and returns the public API router.
func InitPublicAPIRouter(env *config.SEnv, services *services.ServiceComponent) *PublicRouter {
	logger.Debug("Initializing public API router ...")

	publicApiRouter := &PublicRouter{
		Name:             "public API",
		ServiceComponent: services,
	}
	publicApiRouter.Init()
	publicApiRouter.registerMiddlewares()
	publicApiRouter.registerRoutes()

	logger.Debug("Public API registration complete.")
	return publicApiRouter
}

func (r *PublicRouter) registerMiddlewares() {
	middlewaresList := []echo.MiddlewareFunc{
		middlewares.SlashesMiddleware(),
		middlewares.LoggerMiddleware(),
		middlewares.TimeoutMiddleware(),
		middlewares.RequestHeadersMiddleware(),
		middlewares.ResponseHeadersMiddleware(),
	}

	for _, middleware := range middlewaresList {
		r.RegisterMiddleware(middleware)
	}

	if config.DevModeFlag {
		r.RegisterMiddleware(middlewares.BodyDumpMiddleware())
	}
}

// Consolidate route registrations
func (r *PublicRouter) registerRoutes() {
	r.registerPublicApiHealthCheckHandlers()
	r.registerPublicAPIRoutes()
	r.registerPublicApiErrorHandlers()
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
	satelliteHandler := satellites.NewSatelliteHandler(r.ServiceComponent.SatelliteService)
	contextHandler := apicontext.NewContextHandler(r.ServiceComponent.ContextService)
	tileHandler := tiles.NewTileHandler(r.ServiceComponent.TileService)

	satellite := r.Echo.Group("/satellites")
	satellite.GET("/orbit", satelliteHandler.GetSatellitePositionsByNoradID)
	satellite.GET("/paginated", satelliteHandler.GetPaginatedSatellites)
	satellite.GET("/paginated/tles", satelliteHandler.GetPaginatedSatelliteInfo)

	tile := r.Echo.Group("/tiles")
	tile.GET("/all", tileHandler.GetAllTiles)
	tile.GET("/region", tileHandler.GetTilesInRegionHandler)
	tile.GET("/mappings", tileHandler.GetPaginatedSatelliteMappings)
	tile.PUT("/mappings/recompute/bynoradID", tileHandler.RecomputeMappingsByNoradID)
	tile.GET("/mappings/bynoradID", tileHandler.GetSatelliteMappingsByNoradID)

	context := r.Echo.Group("/contexts")
	context.GET("/all", contextHandler.GetPaginatedContexts)
	context.POST("/", contextHandler.CreateContext)
	context.PUT("/:name", contextHandler.UpdateContext)
	context.GET("/:name", contextHandler.GetContextByName)
	context.DELETE("/:name", contextHandler.DeleteContextByName)
	context.PUT("/:name/activate", contextHandler.ActivateContext)
	context.PUT("/:name/deactivate", contextHandler.DeactivateContext)
}

// Middleware helpers
func (r *PublicRouter) RegisterPreMiddleware(middleware echo.MiddlewareFunc) {
	r.Echo.Pre(middleware)
}

// RegisterMiddleware registers middleware
func (r *PublicRouter) RegisterMiddleware(middleware echo.MiddlewareFunc) {
	r.Echo.Use(middleware)
}

// Start the Echo server
func (r *PublicRouter) Start(host string, port string) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		r.Echo.Logger.Info(fmt.Sprintf("Starting %s server on port: %s", r.Name, port))
		if err := r.Echo.Start(host + ":" + port); err != nil && err != http.ErrServerClosed {
			r.Echo.Logger.Fatalf("Server error: %v", err)
		}
	}()

	<-ctx.Done()
	r.shutdownServer()
}

func (r *PublicRouter) shutdownServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := r.Echo.Shutdown(ctx); err != nil {
		r.Echo.Logger.Fatalf("Failed to shutdown server: %v", err)
	}

	r.Echo.Logger.Info("Server shutdown complete.")
}
