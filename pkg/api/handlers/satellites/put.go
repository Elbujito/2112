package satellites

import (
	"net/http"

	"github.com/Elbujito/2112/pkg/api/handlers"
	"github.com/Elbujito/2112/pkg/api/helpers"
	"github.com/Elbujito/2112/pkg/db/models"
	"github.com/Elbujito/2112/pkg/utils/constants"
	"github.com/labstack/echo/v4"
)

// Put handles updating an existing satellite
func Put(c echo.Context) error {
	// Retrieve satellite ID from the request
	id := c.Param("id")
	if id == "" {
		return helpers.Error(c, constants.ERROR_ID_NOT_FOUND, nil)
	}

	// Parse and validate the request body
	form := &models.SatelliteForm{}
	if err := c.Bind(form); err != nil {
		return helpers.Error(c, constants.ERROR_BINDING_BODY, err)
	}

	if err := helpers.Validate(form); err != nil {
		return c.JSON(http.StatusBadRequest, handlers.ValidationErrors(err))
	}

	// Use SatelliteService to find the existing satellite
	service := models.SatelliteModel()
	satellite, err := service.Find(id)
	if err != nil {
		return helpers.Error(c, constants.ERROR_ID_NOT_FOUND, err)
	}

	// Update the satellite with data from the form
	satellite.Name = form.Name
	satellite.NoradID = form.NoradID
	satellite.Type = form.Type

	// Save the updated satellite
	if err := service.Update(satellite); err != nil {
		return helpers.Error(c, constants.ERROR_ID_NOT_FOUND, err)
	}

	// Return success response
	return c.JSON(http.StatusOK, handlers.Success(satellite.MapToForm()))
}
