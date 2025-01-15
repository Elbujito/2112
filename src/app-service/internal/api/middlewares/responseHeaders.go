package middlewares

import (
	xconstants "github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xconstants"
	"github.com/labstack/echo/v4"
)

// ResponseHeadersMiddleware returns Response Middleware
func ResponseHeadersMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(xconstants.HEADER_CONTENT_TYPE, xconstants.HEADER_CONTENT_TYPE_JSON)
			return next(c)
		}
	}
}
