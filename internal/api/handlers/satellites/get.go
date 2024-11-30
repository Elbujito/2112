package tiles

import (
	"net/http"
	"time"

	"github.com/Elbujito/2112/internal/api/handlers"
	"github.com/Elbujito/2112/internal/data"
	repository "github.com/Elbujito/2112/internal/repositories"
	"github.com/Elbujito/2112/internal/services"
	"github.com/Elbujito/2112/pkg/fx/constants"
	"github.com/labstack/echo/v4"
)

// GetSatellitePositionsByNoradID fetches tiles for a given NORAD ID
func GetSatellitePositionsByNoradID(c echo.Context) error {
	// Extract NORAD ID parameter
	noradID := c.QueryParam("noradID")
	if noradID == "" {
		c.Echo().Logger.Error(constants.ERROR_ID_NOT_FOUND)
		return constants.ERROR_ID_NOT_FOUND
	}

	// Assuming you have a service or repository to fetch tiles by NORAD ID
	database := data.NewDatabase()
	tleRepo := repository.NewTLERepository(&database)
	satService := services.NewSatelliteService(tleRepo)
	positions, err := satService.Propagate(c.Request().Context(), noradID, 1*time.Hour)
	if err != nil {
		c.Echo().Logger.Error("Failed to propagate positions: ", err)
		return err
	}

	// If no positions are found
	if len(positions) == 0 {
		c.Echo().Logger.Error(constants.ERROR_ID_NOT_FOUND)
		return constants.ERROR_ID_NOT_FOUND
	}

	// Return positions in the response
	return c.JSON(http.StatusOK, handlers.Success(positions))
}
