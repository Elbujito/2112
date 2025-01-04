package middlewares

import (
	"net/http"
	"strings"

	"github.com/Elbujito/2112/src/templates/go-server/internal/clients/cors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CORSMiddleware() echo.MiddlewareFunc {
	corsCli := cors.GetClient()
	config := corsCli.GetConfig()

	// Split AllowOrigins into a slice if it's a comma-separated string
	var allowedOrigins []string
	if config.AllowOrigins != "" {
		allowedOrigins = strings.Split(config.AllowOrigins, ",")
	}

	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"Content-Type", "Authorization", "Cookie"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Type", "Set-Cookie"},
		MaxAge:           3600, // Cache preflight response for 1 hour
	})
}
