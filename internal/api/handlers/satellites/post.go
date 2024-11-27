package satellites

import (
	"net/http"

	"github.com/Elbujito/2112/internal/api/handlers"
	"github.com/Elbujito/2112/internal/api/helpers"
	"github.com/Elbujito/2112/internal/data/models"
	"github.com/Elbujito/2112/pkg/fx/constants"
	"github.com/labstack/echo/v4"
)

// Post handles the creation of a new satellite
func Post(c echo.Context) error {
	// Parse and validate the request body
	form := &models.SatelliteForm{}
	if err := c.Bind(form); err != nil {
		return helpers.Error(c, constants.ERROR_BINDING_BODY, err)
	}

	if err := helpers.Validate(form); err != nil {
		return c.JSON(http.StatusBadRequest, handlers.ValidationErrors(err))
	}

	// Map the form to the model
	model := form.MapToModel()

	// Use the SatelliteService to save the satellite
	service := models.SatelliteModel()
	if err := service.Save(model); err != nil {
		return helpers.Error(c, constants.ERROR_ID_NOT_FOUND, err)
	}

	// Return success response
	return c.JSON(http.StatusOK, handlers.Success(model.MapToForm()))
}
