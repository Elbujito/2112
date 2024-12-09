package routers

import (
	"github.com/Elbujito/2112/src/template/go-server/internal/api/handlers/errors"
	healthHandlers "github.com/Elbujito/2112/src/template/go-server/internal/api/handlers/healthz"
	"github.com/Elbujito/2112/src/template/go-server/internal/api/middlewares"
	"github.com/Elbujito/2112/src/template/go-server/internal/clients/logger"
	"github.com/Elbujito/2112/src/template/go-server/internal/config"
	"github.com/Elbujito/2112/src/template/go-server/pkg/fx/xconstants"
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

func registerHiddenAPIMiddlewares() {
	hiddenApiRouter.RegisterPreMiddleware(middlewares.SlashesMiddleware())

	hiddenApiRouter.RegisterMiddleware(middlewares.LoggerMiddleware())
	hiddenApiRouter.RegisterMiddleware(middlewares.TimeoutMiddleware())
	hiddenApiRouter.RegisterMiddleware(middlewares.RequestHeadersMiddleware())
	hiddenApiRouter.RegisterMiddleware(middlewares.ResponseHeadersMiddleware())

	if config.Feature(xconstants.FEATURE_GZIP).IsEnabled() {
		hiddenApiRouter.RegisterMiddleware(middlewares.GzipMiddleware())
	}
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

	// if config.Feature(xconstants.FEATURE_ORY_KETO).IsEnabled() {
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

}
