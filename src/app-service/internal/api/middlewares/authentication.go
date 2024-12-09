package middlewares

import (
	"fmt"

	"github.com/Elbujito/2112/src/app-service/internal/clients/kratos"
	xconstants "github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xconstants"

	"github.com/labstack/echo/v4"
)

func AuthenticationMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// skip authentication for health check
			if c.Path() == xconstants.NAME_HEALTH_PATH || c.Path() == fmt.Sprintf("%s%s", xconstants.NAME_HEALTH_PATH, xconstants.NAME_HEALTH_READY_PATH) {
				return next(c)
			}
			// validate session
			kratosCli := kratos.GetClient()
			session, err := kratosCli.ValidateSession(c.Request())
			if err != nil {
				c.Logger().Warn(err)
				c.Logger().Error(xconstants.ERROR_SESSION_NOT_FOUND)
				return xconstants.ERROR_NOT_AUTHORIZED
			}
			if !*session.Active {
				return xconstants.ERROR_NOT_AUTHORIZED
			}
			c.Logger().Warn("Session found:")
			c.Logger().Warn(session)
			kratosCli.Session.SetSession(session)
			return next(c)
		}
	}
}
