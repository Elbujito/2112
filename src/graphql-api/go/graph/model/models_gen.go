// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Mutation struct {
}

type Query struct {
}

type SatellitePosition struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
	Timestamp string  `json:"timestamp"`
}

type SatelliteTle struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	TleLine1 string `json:"tleLine1"`
	TleLine2 string `json:"tleLine2"`
}

type SatelliteVisibility struct {
	SatelliteID   string        `json:"satelliteId"`
	SatelliteName string        `json:"satelliteName"`
	Aos           string        `json:"aos"`
	Los           string        `json:"los"`
	UserLocation  *UserLocation `json:"userLocation"`
}

type Subscription struct {
}

type UserLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    float64 `json:"radius"`
	Horizon   float64 `json:"horizon"`
}

type UserLocationInput struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    float64 `json:"radius"`
	Horizon   float64 `json:"horizon"`
}
