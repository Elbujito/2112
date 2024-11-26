package celestrack

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Elbujito/2112/pkg/api/mappers"
	"github.com/labstack/echo/v4"
)

const CELESTRACK_URL = "https://celestrak.com/NORAD/elements/gp.php"

func FetchTLEHandler(c echo.Context) error {
	noradID := c.Param("norad_id")
	if noradID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "NORAD ID is required"})
	}

	// CelesTrak API endpoint for NORAD-specific TLE
	url := fmt.Sprintf(CELESTRACK_URL+"?CATNR=%s", noradID)

	// Make the API request
	resp, err := http.Get(url)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch TLE data"})
	}
	defer resp.Body.Close()

	// Return the response from CelesTrak
	if resp.StatusCode != http.StatusOK {
		return c.JSON(resp.StatusCode, map[string]string{"error": "Failed to fetch TLE data from CelesTrak"})
	}

	// Stream response directly to the client
	return c.Stream(resp.StatusCode, resp.Header.Get("Content-Type"), resp.Body)
}

func FetchCategoryTLEHandler(c echo.Context) error {
	category := c.Param("category")
	if category == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Category is required"})
	}

	// CelesTrak API endpoint for categorized TLE
	url := fmt.Sprintf(CELESTRACK_URL+"?GROUP=%s", category)

	// Make the API request
	resp, err := http.Get(url)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch TLE data"})
	}
	defer resp.Body.Close()

	// Return the response from CelesTrak
	if resp.StatusCode != http.StatusOK {
		return c.JSON(resp.StatusCode, map[string]string{"error": "Failed to fetch TLE data from CelesTrak"})
	}

	// Stream response directly to the client
	return c.Stream(resp.StatusCode, resp.Header.Get("Content-Type"), resp.Body)
}

// FetchTLE fetches TLE data for a given NORAD ID and returns the parsed TLE domain object.
func FetchTLE(noradID string) (*mappers.RawTLE, error) {
	if noradID == "" {
		return nil, fmt.Errorf("norad_id is required")
	}

	// CelesTrak API endpoint for NORAD-specific TLE
	url := fmt.Sprintf("%s?CATNR=%s", CELESTRACK_URL, noradID)

	// Make the API request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch TLE data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch TLE data: HTTP status %d", resp.StatusCode)
	}

	// Parse TLE response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read TLE data: %v", err)
	}

	lines := bytes.Split(body, []byte("\n"))
	if len(lines) < 2 {
		return nil, fmt.Errorf("TLE data is incomplete")
	}

	line1 := strings.TrimSpace(string(lines[0]))
	line2 := strings.TrimSpace(string(lines[1]))

	return &mappers.RawTLE{
		Line1: line1,
		Line2: line2,
	}, nil
}

// FetchCategoryTLE fetches TLEs for a given category from CelesTrak
func FetchCategoryTLE(category string) ([]*mappers.RawTLE, error) {
	if category == "" {
		return nil, fmt.Errorf("category is required")
	}

	// Construct the URL for the category
	url := fmt.Sprintf("%s?GROUP=%s", CELESTRACK_URL, category)

	// Fetch the TLE data
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch TLE data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch TLE data: HTTP status %d", resp.StatusCode)
	}

	// Parse the TLE data
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read TLE data: %v", err)
	}

	lines := bytes.Split(body, []byte("\n"))
	var tles []*mappers.RawTLE
	for i := 0; i < len(lines)-1; i += 3 {
		if len(lines[i]) == 0 || len(lines[i+1]) == 0 || len(lines[i+2]) == 0 {
			continue
		}
		noradID := strings.TrimSpace(string(lines[i]))
		line1 := strings.TrimSpace(string(lines[i+1]))
		line2 := strings.TrimSpace(string(lines[i+2]))

		tles = append(tles, &mappers.RawTLE{
			NoradID: noradID,
			Line1:   line1,
			Line2:   line2,
		})
	}

	return tles, nil
}
