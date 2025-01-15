package routers

import (
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/api/handlers/errors"
	healthHandlers "github.com/Elbujito/2112/src/app-service/internal/api/handlers/healthz"
	"github.com/Elbujito/2112/src/app-service/internal/api/middlewares"
	"github.com/Elbujito/2112/src/app-service/internal/clients/logger"
	"github.com/Elbujito/2112/src/app-service/internal/config"
	xconstants "github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xconstants"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var hiddenApiRouter *PublicRouter

func InitHiddenAPIRouter() {
	logger.Debug("Initializing hidden api router ...")
	hiddenApiRouter = &PublicRouter{}
	hiddenApiRouter.Name = "hidden API"
	hiddenApiRouter.Init()

	// order is important here
	// first register development middlewares
	if config.DevModeFlag {
		logger.Debug("Registering hidden api development middlewares ...")
		registerHiddenApiDevModeMiddleware()
	}

	// next register middlwares
	logger.Debug("Registering hidden api middlewares ...")
	registerHiddenAPIMiddlewares()

	// next register all health check routes
	logger.Debug("Registering hidden api health routes ...")
	registerHiddenApiHealthCheckHandlers()

	// next register security related middleware
	logger.Debug("Registering hidden api security middlewares ...")
	registerHiddenApiSecurityMiddlewares()

	// next register all routes
	logger.Debug("Registering hidden api hidden routes ...")
	registerHiddenAPIRoutes()

	// finally register default fallback error handlers
	// 404 is handled here as the last route
	logger.Debug("Registering hidden api error handlers ...")
	registerHiddenApiErrorHandlers()

	logger.Debug("Hidden api registration complete.")
}

func HiddenAPIRouter() *PublicRouter {
	return hiddenApiRouter
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

func registerHiddenAPIMiddlewares() {
	hiddenApiRouter.RegisterPreMiddleware(middlewares.SlashesMiddleware())
	hiddenApiRouter.RegisterMiddleware(middlewares.LoggerMiddleware())
	hiddenApiRouter.RegisterMiddleware(middlewares.TimeoutMiddleware())
	hiddenApiRouter.registerPrometheusMetrics()
}

func registerHiddenApiDevModeMiddleware() {
	hiddenApiRouter.RegisterMiddleware(middlewares.BodyDumpMiddleware())
}

func registerHiddenApiSecurityMiddlewares() {
	hiddenApiRouter.RegisterMiddleware(middlewares.XSSCheckMiddleware())

	if config.Feature(xconstants.FEATURE_CORS).IsEnabled() {
		hiddenApiRouter.RegisterMiddleware(middlewares.CORSMiddleware())
	}

	if config.Feature(xconstants.FEATURE_ORY_KRATOS).IsEnabled() {
		hiddenApiRouter.RegisterMiddleware(middlewares.AuthenticationMiddleware())
	}
}

func registerHiddenApiErrorHandlers() {
	hiddenApiRouter.Echo.HTTPErrorHandler = errors.AutomatedHttpErrorHandler()
	hiddenApiRouter.Echo.RouteNotFound("/*", errors.NotFound)
}

func registerHiddenApiHealthCheckHandlers() {
	health := hiddenApiRouter.Echo.Group("/health")
	health.GET("/alive", healthHandlers.Index)
	health.GET("/ready", healthHandlers.Ready)
}

func registerHiddenAPIRoutes() {
	hiddenApiRouter.Echo.GET("/metrics", func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, "text/plain; charset=utf-8")
		promhttp.Handler().ServeHTTP(c.Response(), c.Request())
		return nil
	})
}
