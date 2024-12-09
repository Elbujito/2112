package middlewares

import (
	"github.com/Elbujito/2112/template/go-server/pkg/fx/xconstants"
	"github.com/labstack/echo/v4"
)

func ResponseHeadersMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(xconstants.HEADER_CONTENT_TYPE, xconstants.HEADER_CONTENT_TYPE_JSON)
			return next(c)
		}
	}
}
