package middlewares

import (
	"net/http"
	"strings"

	"github.com/Elbujito/2112/src/app-service/internal/clients/cors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

// CORSMiddleware returns CORS Middleware
func CORSMiddleware() echo.MiddlewareFunc {
	corsCli := cors.GetClient()
	config := corsCli.GetConfig()

	var allowedOrigins []string
	if config.AllowOrigins != "" {
		for _, origin := range strings.Split(config.AllowOrigins, ",") {
			trimmedOrigin := strings.TrimSpace(origin)
			if trimmedOrigin != "" {
				allowedOrigins = append(allowedOrigins, trimmedOrigin)
			}
		}
	}

	if len(allowedOrigins) == 0 {
		logrus.Warn("CORS AllowOrigins is empty. Defaulting to wildcard (*) for development.")
		allowedOrigins = append(allowedOrigins, "*")
	}

	logrus.Infof("CORS Config: Allowed Origins: %v", allowedOrigins)

	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "Authorization", "Cookie"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Type", "Set-Cookie"},
		MaxAge:           3600,
	})
}
