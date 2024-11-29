package tiles

import (
	"net/http"

	"github.com/Elbujito/2112/internal/api/handlers"
	"github.com/Elbujito/2112/internal/data"
	repository "github.com/Elbujito/2112/internal/repositories"
	"github.com/Elbujito/2112/pkg/fx/constants"
	"github.com/labstack/echo/v4"
)

// GetTileTest fetches tiles for a given NORAD ID
func GetTileTest(c echo.Context) error {
	return c.JSON(http.StatusOK, handlers.Success("test"))
}

// GetTilesByNoradID fetches tiles for a given NORAD ID
func GetTilesByNoradID(c echo.Context) error {
	// Extract NORAD ID parameter
	noradID := c.QueryParam("noradID")
	if noradID == "" {
		c.Echo().Logger.Error(constants.ERROR_ID_NOT_FOUND)
		return constants.ERROR_ID_NOT_FOUND
	}

	// Assuming you have a service or repository to fetch tiles by NORAD ID
	database := data.NewDatabase()
	tileSatelliteMappingRepo := repository.NewTileSatelliteMappingRepository(&database)
	tiles, err := tileSatelliteMappingRepo.FindAllVisibleTilesByNoradIDSortedByAOSTime(c.Request().Context(), noradID)
	if err != nil {
		c.Echo().Logger.Error("Failed to fetch tiles: ", err)
		return err
	}

	// If no tiles are found
	if len(tiles) == 0 {
		c.Echo().Logger.Error(constants.ERROR_ID_NOT_FOUND)
		return constants.ERROR_ID_NOT_FOUND
	}

	// Return tiles in the response
	return c.JSON(http.StatusOK, handlers.Success(tiles))
}
