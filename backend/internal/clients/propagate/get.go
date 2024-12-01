package propagator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const DefaultPropagationAPIURL = "http://2112-propagator:5000/satellite/propagate"

type PropagatorClient struct {
	BaseURL string
}

// SatellitePropagationRequest represents the payload for the propagation request.
type SatellitePropagationRequest struct {
	TLELine1        string `json:"tle_line1"`
	TLELine2        string `json:"tle_line2"`
	StartTime       string `json:"start_time"`
	DurationMinutes int    `json:"duration_minutes"`
	IntervalSeconds int    `json:"interval_seconds"`
}

// SatellitePosition represents the propagated position of a satellite.
type SatellitePosition struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
	Time      string  `json:"time"`
}

// NewPropagatorClient creates a new PropagatorClient with the given base URL.
func NewPropagatorClient(baseURL string) *PropagatorClient {
	if baseURL == "" {
		baseURL = DefaultPropagationAPIURL
	}
	return &PropagatorClient{
		BaseURL: baseURL,
	}
}

// FetchPropagation fetches propagated positions for a given TLE and parameters.
func (client *PropagatorClient) FetchPropagation(ctx context.Context, tle1, tle2, startTime string, durationMinutes, intervalSeconds int) ([]*SatellitePosition, error) {
	// Validate input
	if tle1 == "" || tle2 == "" || startTime == "" {
		return nil, fmt.Errorf("TLE lines and start time are required")
	}
	if durationMinutes <= 0 || intervalSeconds <= 0 {
		return nil, fmt.Errorf("duration and interval must be greater than zero")
	}

	// Prepare request payload
	requestPayload := SatellitePropagationRequest{
		TLELine1:        tle1,
		TLELine2:        tle2,
		StartTime:       startTime,
		DurationMinutes: durationMinutes,
		IntervalSeconds: intervalSeconds,
	}
	payloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %v", err)
	}

	// Send the HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, client.BaseURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch propagation data: HTTP status %d, response: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var positions []*SatellitePosition
	if err := json.NewDecoder(resp.Body).Decode(&positions); err != nil {
		return nil, fmt.Errorf("failed to parse response body: %v", err)
	}

	return positions, nil
}
