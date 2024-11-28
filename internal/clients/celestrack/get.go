package celestrack

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Elbujito/2112/internal/api/mappers"
)

const CELESTRACK_URL = "https://celestrak.com/NORAD/elements/gp.php"

type CelestrackClient struct {
}

// FetchTLEFromSatCatByCategory fetches TLEs for a given category from CelesTrak
func (client *CelestrackClient) FetchTLEFromSatCatByCategory(ctx context.Context, category string) ([]*mappers.RawTLE, error) {
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
