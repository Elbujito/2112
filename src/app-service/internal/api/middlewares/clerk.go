package middlewares

import (
	"net/http"

	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/labstack/echo/v4"
)

// ClerkMiddleware is a wrapper to adapt Clerk middleware to Echo's middleware interface.
func ClerkMiddleware(opts ...clerkhttp.AuthorizationOption) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			handler := clerkhttp.RequireHeaderAuthorization(opts...)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.SetRequest(r)
				c.Response().Writer = w
				_ = next(c)
			}))

			handler.ServeHTTP(c.Response().Writer, c.Request())
			return nil
		}
	}
}
