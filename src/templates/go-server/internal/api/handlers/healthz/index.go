package healthz

import (
	"net/http"

	"github.com/Elbujito/2112/src/template/go-server/internal/api/handlers"
	"github.com/Elbujito/2112/src/template/go-server/internal/config"

	"github.com/labstack/echo/v4"
)

func Index(c echo.Context) error {
	payload := map[string]string{
		"message": "ok",
		"version": config.Env.Version,
	}

	return c.JSON(http.StatusOK, handlers.Success(payload))
}
