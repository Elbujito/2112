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
	serviceapi "github.com/Elbujito/2112/src/app-service/internal/api/services"
	"github.com/Elbujito/2112/src/app-service/internal/clients/logger"
	"github.com/Elbujito/2112/src/app-service/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PublicRouter manages the public API router and its dependencies.
type PublicRouter struct {
	Echo             *echo.Echo
	Name             string
	ServiceComponent *serviceapi.ServiceComponent
}

// Init initializes the Echo instance for the router.
func (r *PublicRouter) Init() {
	r.Echo = echo.New()
	r.Echo.HideBanner = true
	r.Echo.Logger = logger.GetLogger()

}

func (r *PublicRouter) registerPrometheusMetrics() {
	// Metrics to monitor the application and database health
	httpRequestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests processed, labeled by status code and method.",
		},
		[]string{"method", "endpoint", "status"},
	)
	httpRequestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status"},
	)
	httpRequestSize := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP request size in bytes.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
	httpResponseSize := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status"},
	)
	appUptime := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "app_uptime_seconds",
			Help: "Application uptime in seconds.",
		},
	)
	appLiveness := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "app_liveness",
			Help: "Application liveness status (1 = alive, 0 = not alive).",
		},
	)

	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(httpRequestSize)
	prometheus.MustRegister(httpResponseSize)
	prometheus.MustRegister(appUptime)
	prometheus.MustRegister(appLiveness)

	go func() {
		startTime := time.Now()
		appLiveness.Set(1) // Set liveness to 1 indicating the app is alive
		for {
			appUptime.Set(time.Since(startTime).Seconds())
			time.Sleep(1 * time.Second)
		}
	}()
}

// InitPublicAPIRouter initializes and returns the public API router.
func InitPublicAPIRouter(env *config.SEnv) *PublicRouter {
	logger.Debug("Initializing public API router ...")

	serviceComponent := serviceapi.NewServiceComponent(env)
	publicApiRouter := &PublicRouter{
		Name:             "public API",
		ServiceComponent: serviceComponent,
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
		// middlewares.RequestHeadersMiddleware(),
		// middlewares.ResponseHeadersMiddleware(),
	}

	// if config.Feature(xconstants.FEATURE_GZIP).IsEnabled() {
	// 	middlewaresList = append(middlewaresList, middlewares.GzipMiddleware())
	// }

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
	r.registerPrometheusMetrics()
	r.registerMetricsEndpoint()
	r.registerPublicAPIRoutes()
	r.registerPublicApiErrorHandlers()
}

// Register health check handlers
func (r *PublicRouter) registerPublicApiHealthCheckHandlers() {
	health := r.Echo.Group("/health")
	health.GET("/alive", healthHandlers.Index)
	health.GET("/ready", healthHandlers.Ready)
}

// Register metrics endpoint
func (r *PublicRouter) registerMetricsEndpoint() {
	r.Echo.GET("/metrics", func(c echo.Context) error {
		// Set the Content-Type to text/plain; charset=utf-8
		c.Response().Header().Set(echo.HeaderContentType, "text/plain; charset=utf-8")
		// Use promhttp.Handler to serve the metrics
		promhttp.Handler().ServeHTTP(c.Response(), c.Request())
		return nil
	})
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
