package users

import (
	"net/http"

	"github.com/Elbujito/2112/src/templates/go-server/internal/api/handlers"
	"github.com/Elbujito/2112/src/templates/go-server/internal/clients/kratos"

	"github.com/labstack/echo/v4"
)

func Index(c echo.Context) error {
	kratosCli := kratos.GetClient()
	identities, err := kratosCli.GetAllIdentity()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, handlers.Success(identities))
}

// Delete handler
