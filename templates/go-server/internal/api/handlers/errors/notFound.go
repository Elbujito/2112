package errors

import (
	"net/http"

	"github.com/Elbujito/2112/fx/constants"
	"github.com/Elbujito/2112/template/go-server/internal/api/handlers"

	"github.com/labstack/echo/v4"
)

func NotFound(c echo.Context) error {
	return c.JSON(
		http.StatusNotFound,
		handlers.BuildResponse(
			constants.STATUS_CODE_ROUTE_NOT_FOUND,
			constants.MSG_ROUTE_NOT_FOUND,
			[]string{},
			nil))
}
