package middlewares

import (
	"github.com/Elbujito/2112/pkg/clients/gzip"
	"github.com/Elbujito/2112/pkg/fx"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func GzipMiddleware() echo.MiddlewareFunc {
	GzipCli := gzip.GetClient()
	config := GzipCli.GetConfig()
	level := fx.IntFromStr(config.Level)
	return middleware.GzipWithConfig(middleware.GzipConfig{Level: level})
}
