package satellites

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/domain"
	"github.com/Elbujito/2112/src/app-service/internal/services"
	xconstants "github.com/Elbujito/2112/src/templates/go-server/pkg/fx/xconstants"
	"github.com/labstack/echo/v4"
)

type SatelliteHandler struct {
	Service services.SatelliteService
}

// NewSatelliteHandler creates a new handler with the provided SatelliteService.
func NewSatelliteHandler(service services.SatelliteService) *SatelliteHandler {
	return &SatelliteHandler{Service: service}
}

// GetSatellitePositionsByNoradID fetches satellite positions by NORAD ID.
func (h *SatelliteHandler) GetSatellitePositionsByNoradID(c echo.Context) error {
	noradID := c.QueryParam("noradID")
	if noradID == "" {
		c.Echo().Logger.Error(xconstants.ERROR_ID_NOT_FOUND)
		return xconstants.ERROR_ID_NOT_FOUND
	}

	positions, err := h.Service.Propagate(c.Request().Context(), noradID, 24*time.Hour, 1*time.Minute)
	if err != nil {
		c.Echo().Logger.Error("Failed to propagate positions: ", err)
		return err
	}

	if len(positions) == 0 {
		c.Echo().Logger.Error(xconstants.ERROR_ID_NOT_FOUND)
		return xconstants.ERROR_ID_NOT_FOUND
	}

	return c.JSON(http.StatusOK, positions)
}

// GetPaginatedSatellites fetches a paginated list of satellites with optional search filters.
func (h *SatelliteHandler) GetPaginatedSatellites(c echo.Context) error {
	// Parse query parameters for pagination
	pageStr := c.QueryParam("page")
	pageSizeStr := c.QueryParam("pageSize")
	searchWildcard := c.QueryParam("search") // Retrieve optional search query

	// Convert parameters to integers with defaults
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1 // Default to page 1 if invalid
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10 // Default to 10 records per page if invalid
	}

	// Create SearchRequest object
	searchRequest := &domain.SearchRequest{
		Wildcard: searchWildcard,
	}

	// Call the service method for pagination with search filters
	satellites, totalRecords, err := h.Service.ListSatellitesWithPagination(c.Request().Context(), page, pageSize, searchRequest)
	if err != nil {
		c.Echo().Logger.Error("Failed to fetch paginated satellites: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to fetch satellites")
	}

	// Prepare the response
	response := map[string]interface{}{
		"totalRecords": totalRecords,
		"page":         page,
		"pageSize":     pageSize,
		"satellites":   satellites,
	}

	return c.JSON(http.StatusOK, response)
}
