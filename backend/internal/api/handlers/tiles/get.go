package tiles

import (
	"net/http"
	"strconv"

	"github.com/Elbujito/2112/internal/services"
	"github.com/Elbujito/2112/lib/fx/constants"
	"github.com/labstack/echo/v4"
)

type TileHandler struct {
	Service services.TileService
}

// NewTileHandler creates a new handler with the provided TileService.
func NewTileHandler(service services.TileService) *TileHandler {
	return &TileHandler{Service: service}
}

// GetAllTiles fetches all available tiles.
func (h *TileHandler) GetAllTiles(c echo.Context) error {
	tiles, err := h.Service.FindAllTiles(c.Request().Context())
	if err != nil {
		c.Echo().Logger.Error("Failed to fetch tiles: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to fetch tiles")
	}

	// If no tiles are found
	if len(tiles) == 0 {
		c.Echo().Logger.Error(constants.ERROR_ID_NOT_FOUND)
		return constants.ERROR_ID_NOT_FOUND
	}

	// Return tiles in the response
	return c.JSON(http.StatusOK, tiles)
}

// GetTilesInRegionHandler handles requests to fetch tiles in a region.
func (h *TileHandler) GetTilesInRegionHandler(c echo.Context) error {
	// Parse query parameters for bounding box
	minLatStr := c.QueryParam("minLat")
	minLonStr := c.QueryParam("minLon")
	maxLatStr := c.QueryParam("maxLat")
	maxLonStr := c.QueryParam("maxLon")

	// Convert query parameters to float64
	minLat, err := strconv.ParseFloat(minLatStr, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid minLat parameter")
	}
	minLon, err := strconv.ParseFloat(minLonStr, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid minLon parameter")
	}
	maxLat, err := strconv.ParseFloat(maxLatStr, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid maxLat parameter")
	}
	maxLon, err := strconv.ParseFloat(maxLonStr, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid maxLon parameter")
	}

	// Call the service to fetch tiles
	tiles, err := h.Service.GetTilesInRegion(c.Request().Context(), minLat, minLon, maxLat, maxLon)
	if err != nil {
		c.Logger().Error("Failed to fetch tiles in region:", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to fetch tiles in region")
	}

	// Return tiles in JSON response
	return c.JSON(http.StatusOK, tiles)
}
