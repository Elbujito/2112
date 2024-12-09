package middlewares

import (
	"github.com/Elbujito/2112/src/app-service/internal/clients/gzip"
	"github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xutils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func GzipMiddleware() echo.MiddlewareFunc {
	GzipCli := gzip.GetClient()
	config := GzipCli.GetConfig()
	level := xutils.IntFromStr(config.Level)
	return middleware.GzipWithConfig(middleware.GzipConfig{Level: level})
}
