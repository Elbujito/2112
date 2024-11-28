package routers

import (
	celestrackHandlers "github.com/Elbujito/2112/internal/api/clients/celestrack"
	"github.com/Elbujito/2112/internal/api/handlers/errors"
	healthHandlers "github.com/Elbujito/2112/internal/api/handlers/healthz"
	"github.com/Elbujito/2112/internal/api/middlewares"
	"github.com/Elbujito/2112/internal/clients/logger"
	"github.com/Elbujito/2112/internal/config"
	"github.com/Elbujito/2112/pkg/fx/constants"
)

var hiddenApiRouter *Router

func InitHiddenAPIRouter() {
	logger.Debug("Initializing hidden api router ...")
	hiddenApiRouter = &Router{}
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

func HiddenAPIRouter() *Router {
	return hiddenApiRouter
}

func registerHiddenAPIMiddlewares() {
	hiddenApiRouter.RegisterPreMiddleware(middlewares.SlashesMiddleware())

	hiddenApiRouter.RegisterMiddleware(middlewares.LoggerMiddleware())
	hiddenApiRouter.RegisterMiddleware(middlewares.TimeoutMiddleware())
	hiddenApiRouter.RegisterMiddleware(middlewares.RequestHeadersMiddleware())
	hiddenApiRouter.RegisterMiddleware(middlewares.ResponseHeadersMiddleware())

	if config.Feature(constants.FEATURE_GZIP).IsEnabled() {
		hiddenApiRouter.RegisterMiddleware(middlewares.GzipMiddleware())
	}
}

func registerHiddenApiDevModeMiddleware() {
	hiddenApiRouter.RegisterMiddleware(middlewares.BodyDumpMiddleware())
}

func registerHiddenApiSecurityMiddlewares() {
	hiddenApiRouter.RegisterMiddleware(middlewares.XSSCheckMiddleware())

	if config.Feature(constants.FEATURE_CORS).IsEnabled() {
		hiddenApiRouter.RegisterMiddleware(middlewares.CORSMiddleware())
	}

	if config.Feature(constants.FEATURE_ORY_KRATOS).IsEnabled() {
		hiddenApiRouter.RegisterMiddleware(middlewares.AuthenticationMiddleware())
	}

	// if config.Feature(constants.FEATURE_ORY_KETO).IsEnabled() {
	// 	// keto middleware <- this will check if the user has the right permissions like system admin
	// 	// hiddenApiRouter.RegisterMiddleware(middlewares.AuthenticationMiddleware())
	// }
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
	// Create a group for CelesTrak-related routes
	celestrak := hiddenApiRouter.Echo.Group("/celestrak")

	// Route to fetch TLE data by NORAD ID
	celestrak.GET("/tle/:norad_id", celestrackHandlers.FetchTLEHandler)

	// Route to fetch categorized TLE data
	celestrak.GET("/categories/:category", celestrackHandlers.FetchCategoryTLEHandler)
}
