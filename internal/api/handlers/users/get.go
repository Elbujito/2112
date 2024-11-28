package users

import (
	"net/http"

	"github.com/Elbujito/2112/internal/api/handlers"
	"github.com/Elbujito/2112/internal/clients/kratos"
	"github.com/Elbujito/2112/pkg/fx/constants"

	"github.com/labstack/echo/v4"
)

func Get(c echo.Context) error {
	id, err := handlers.GetUUIDParam(c.Param("id"))
	if err != nil {
		c.Echo().Logger.Error(constants.ERROR_ID_NOT_FOUND)
		return constants.ERROR_ID_NOT_FOUND
	}
	kratosCli := kratos.GetClient()
	identity, err := kratosCli.GetIdentity(id.String())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, handlers.Success(identity))
}
