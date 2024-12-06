package middlewares

import (
	"time"

	"github.com/Elbujito/2112/template/go-server/internal/clients/service"
	"github.com/Elbujito/2112/template/go-server/pkg/fx"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func TimeoutMiddleware() echo.MiddlewareFunc {
	serviceCli := service.GetClient()
	config := serviceCli.GetConfig()
	timeoutDuration := fx.IntFromStr(config.RequestTimeoutDuration)

	return middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: time.Duration(timeoutDuration) * time.Second,
	})
}
