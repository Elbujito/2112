package satellites

import (
	"net/http"

	"github.com/Elbujito/2112/pkg/api/handlers"
	"github.com/Elbujito/2112/pkg/api/helpers"
	"github.com/Elbujito/2112/pkg/db/models"
	"github.com/Elbujito/2112/pkg/utils/constants"

	"github.com/labstack/echo/v4"
)

func Get(c echo.Context) error {

	id := c.Param("id")

	if id == "" {
		return helpers.Error(c, constants.ERROR_ID_NOT_FOUND, nil)
	}

	m, err := models.SatelliteModel().Find(id)

	if err != nil {
		return helpers.Error(c, err, nil)
	}

	return c.JSON(http.StatusOK, handlers.Success(m.MapToForm()))

}
