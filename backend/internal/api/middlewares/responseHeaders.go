package middlewares

import (
	"github.com/Elbujito/2112/lib/fx/constants"

	"github.com/labstack/echo/v4"
)

func ResponseHeadersMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(constants.HEADER_CONTENT_TYPE, constants.HEADER_CONTENT_TYPE_JSON)
			return next(c)
		}
	}
}
