package routers

import (
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/api/handlers/errors"
	healthHandlers "github.com/Elbujito/2112/src/app-service/internal/api/handlers/healthz"
	metricsHandlers "github.com/Elbujito/2112/src/app-service/internal/api/handlers/metrics"
	"github.com/Elbujito/2112/src/app-service/internal/api/middlewares"
	"github.com/Elbujito/2112/src/app-service/internal/clients/logger"
	"github.com/Elbujito/2112/src/app-service/internal/config"
	xconstants "github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xconstants"
	"github.com/prometheus/client_golang/prometheus"
)

var protectedApiRouter *PublicRouter

func InitProtectedAPIRouter() {
	logger.Debug("Initializing protected api router ...")
	protectedApiRouter = &PublicRouter{}
	protectedApiRouter.Name = "protected API"
	protectedApiRouter.Init()

	// order is important here
	// first register development middlewares
	if config.DevModeFlag {
		logger.Debug("Registering protected api development middlewares ...")
		registerProtectedApiDevModeMiddleware()
	}

	// next register middlwares
	logger.Debug("Registering protected api middlewares ...")
	registerProtectedAPIMiddlewares()

	// next register all health check routes
	logger.Debug("Registering protected api health routes ...")
	registerProtectedApiHealthCheckHandlers()

	// next register security related middleware
	logger.Debug("Registering protected api security middlewares ...")
	registerProtectedApiSecurityMiddlewares()

	// next register all routes
	logger.Debug("Registering protected api protected routes ...")
	registerProtectedAPIRoutes()
	logger.Debug("Registering metrics api protected routes ...")
	registerMetricsAPIRoutes()

	// finally register default fallback error handlers
	// 404 is handled here as the last route
	logger.Debug("Registering protected api error handlers ...")
	registerProtectedApiErrorHandlers()

	logger.Debug("Protected api registration complete.")
}

func ProtectedAPIRouter() *PublicRouter {
	return protectedApiRouter
}

func (r *PublicRouter) registerPrometheusMetrics() {
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

	prometheus.MustRegister(appUptime)
	prometheus.MustRegister(appLiveness)

	go func() {
		startTime := time.Now()
		appLiveness.Set(1)
		for {
			appUptime.Set(time.Since(startTime).Seconds())
			time.Sleep(1 * time.Second)
		}
	}()
}

func registerProtectedAPIMiddlewares() {
	protectedApiRouter.RegisterPreMiddleware(middlewares.SlashesMiddleware())

	protectedApiRouter.RegisterMiddleware(middlewares.LoggerMiddleware())
	protectedApiRouter.RegisterMiddleware(middlewares.TimeoutMiddleware())
	hiddenApiRouter.registerPrometheusMetrics()
}

func registerProtectedApiDevModeMiddleware() {
	protectedApiRouter.RegisterMiddleware(middlewares.BodyDumpMiddleware())
}

func registerProtectedApiSecurityMiddlewares() {
	protectedApiRouter.RegisterMiddleware(middlewares.XSSCheckMiddleware())

	if config.Feature(xconstants.FEATURE_CORS).IsEnabled() {
		protectedApiRouter.RegisterMiddleware(middlewares.CORSMiddleware())
	}
}

func registerProtectedApiErrorHandlers() {
	protectedApiRouter.Echo.HTTPErrorHandler = errors.AutomatedHttpErrorHandler()
	protectedApiRouter.Echo.RouteNotFound("/*", errors.NotFound)
}

func registerProtectedApiHealthCheckHandlers() {
	health := protectedApiRouter.Echo.Group("/health")
	health.GET("/alive", healthHandlers.Index)
	health.GET("/ready", healthHandlers.Ready)
}

func registerProtectedAPIRoutes() {
}

func registerMetricsAPIRoutes() {
	metrics := protectedApiRouter.Echo.Group("/metrics")
	metrics.GET("", metricsHandlers.GetMetrics)
}
