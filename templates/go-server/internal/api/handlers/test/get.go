package test

import (
	"net/http"

	"github.com/Elbujito/2112/template/go-server/internal/services"
	"github.com/Elbujito/2112/template/go-server/pkg/fx/constants"
	"github.com/labstack/echo/v4"
)

type TestHandler struct {
	Service services.TestService
}

// NewSatelliteHandler creates a new handler with the provided TestService.
func NewTestHandler(service services.TestService) *TestHandler {
	return &TestHandler{Service: service}
}

// GetTestByTest fetches tests
func (h *TestHandler) GetTestByTest(c echo.Context) error {
	noradID := c.QueryParam("test")
	if noradID == "" {
		c.Echo().Logger.Error(constants.ERROR_ID_NOT_FOUND)
		return constants.ERROR_ID_NOT_FOUND
	}

	return c.JSON(http.StatusOK, nil)
}
