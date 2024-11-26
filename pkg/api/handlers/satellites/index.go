package cats

import (
	"net/http"

	"github.com/Elbujito/2112/pkg/api/handlers"
	"github.com/Elbujito/2112/pkg/api/helpers"
	"github.com/Elbujito/2112/pkg/db/models"

	"github.com/labstack/echo/v4"
)

func Index(c echo.Context) error {

	ms, err := models.SatelliteModel().FindAll()

	if err != nil {
		return helpers.Error(c, err, nil)
	}

	var payload []*models.SatelliteForm

	for _, m := range ms {
		f := m.MapToForm()
		payload = append(payload, f)
	}

	return c.JSON(http.StatusOK, handlers.Success(payload))

}
