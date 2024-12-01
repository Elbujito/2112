package routers

import (
	"github.com/Elbujito/2112/internal/api/handlers/errors"
	healthHandlers "github.com/Elbujito/2112/internal/api/handlers/healthz"
	satellitesHandlers "github.com/Elbujito/2112/internal/api/handlers/satellites"
	tilesHandlers "github.com/Elbujito/2112/internal/api/handlers/tiles"
	"github.com/Elbujito/2112/internal/api/middlewares"
	"github.com/Elbujito/2112/internal/clients/logger"
	"github.com/Elbujito/2112/internal/config"
	"github.com/Elbujito/2112/pkg/fx/constants"
)

var publicApiRouter *Router

func InitPublicAPIRouter() {
	logger.Debug("Initializing public api router ...")
	publicApiRouter = &Router{}
	publicApiRouter.Name = "public API"
	publicApiRouter.Init()

	// order is important here
	// first register development middlewares
	if config.DevModeFlag {
		logger.Debug("Registering public api development middlewares ...")
		registerPublicApiDevModeMiddleware()
	}

	// next register middlwares
	logger.Debug("Registering public api middlewares ...")
	registerPublicAPIMiddlewares()

	// next register all health check routes
	logger.Debug("Registering public api health routes ...")
	registerPublicApiHealthCheckHandlers()

	// next register security related middleware
	logger.Debug("Registering public api security middlewares ...")
	registerPublicApiSecurityMiddlewares()

	// next register all routes
	logger.Debug("Registering public api public routes ...")
	registerPublicAPIRoutes()

	logger.Debug("Registering public celestrack api handlers ...")
	registerPublicCelestrackAPIRoutes()

	// finally register default fallback error handlers
	// 404 is handled here as the last route
	logger.Debug("Registering public api error handlers ...")
	registerPublicApiErrorHandlers()

	logger.Debug("Public api registration complete.")

}

func PublicAPIRouter() *Router {
	return publicApiRouter
}

func registerPublicAPIMiddlewares() {
	publicApiRouter.RegisterPreMiddleware(middlewares.SlashesMiddleware())

	publicApiRouter.RegisterMiddleware(middlewares.LoggerMiddleware())
	publicApiRouter.RegisterMiddleware(middlewares.TimeoutMiddleware())
	publicApiRouter.RegisterMiddleware(middlewares.RequestHeadersMiddleware())
	publicApiRouter.RegisterMiddleware(middlewares.ResponseHeadersMiddleware())

	if config.Feature(constants.FEATURE_GZIP).IsEnabled() {
		publicApiRouter.RegisterMiddleware(middlewares.GzipMiddleware())
	}
}

func registerPublicApiDevModeMiddleware() {
	publicApiRouter.RegisterMiddleware(middlewares.BodyDumpMiddleware())
}

func registerPublicApiSecurityMiddlewares() {
	publicApiRouter.RegisterMiddleware(middlewares.XSSCheckMiddleware())

	if config.Feature(constants.FEATURE_CORS).IsEnabled() {
		publicApiRouter.RegisterMiddleware(middlewares.CORSMiddleware())
	}

}

func registerPublicApiErrorHandlers() {
	publicApiRouter.Echo.HTTPErrorHandler = errors.AutomatedHttpErrorHandler()
	publicApiRouter.Echo.RouteNotFound("/*", errors.NotFound)
}

func registerPublicApiHealthCheckHandlers() {
	health := publicApiRouter.Echo.Group("/health")
	health.GET("/alive", healthHandlers.Index)
	health.GET("/ready", healthHandlers.Ready)
}

func registerPublicAPIRoutes() {
	tile := publicApiRouter.Echo.Group("/tiles")
	tile.GET("/mapping", tilesHandlers.GetTilesByNoradID)
	tile.GET("/all", tilesHandlers.GetTiles)
	satellite := publicApiRouter.Echo.Group("/satellites")
	satellite.GET("/orbit", satellitesHandlers.GetSatellitePositionsByNoradID)
}

func registerPublicCelestrackAPIRoutes() {

}
