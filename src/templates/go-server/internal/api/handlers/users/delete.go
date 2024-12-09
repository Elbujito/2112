package users

import (
	"net/http"

	"github.com/Elbujito/2112/src/templates/go-server/internal/api/handlers"
	"github.com/Elbujito/2112/src/templates/go-server/internal/clients/kratos"
	"github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xconstants"

	"github.com/labstack/echo/v4"
)

func Delete(c echo.Context) error {
	id, err := handlers.GetUUIDParam(c.Param("id"))
	if err != nil {
		c.Logger().Error(xconstants.ERROR_ID_NOT_FOUND)
		return xconstants.ERROR_ID_NOT_FOUND
	}
	kratosCli := kratos.GetClient()
	if err := kratosCli.DeleteIdentity(id.String()); err != nil {
		return err
	}
	return c.JSON(http.StatusAccepted, handlers.Accepted())
}
