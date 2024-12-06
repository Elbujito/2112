package middlewares

import (
	"github.com/Elbujito/2112/lib/fx/xutils"
	"github.com/Elbujito/2112/template/go-server/internal/clients/gzip"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func GzipMiddleware() echo.MiddlewareFunc {
	GzipCli := gzip.GetClient()
	config := GzipCli.GetConfig()
	level := xutils.IntFromStr(config.Level)
	return middleware.GzipWithConfig(middleware.GzipConfig{Level: level})
}
