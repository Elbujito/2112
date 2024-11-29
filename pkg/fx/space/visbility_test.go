package space

import (
	"math"
	"testing"
	"time"

	"github.com/Elbujito/2112/pkg/fx/polygon"
	"github.com/joshuaferrara/go-satellite"
)

func TestHaversineDistance(t *testing.T) {
	tests := []struct {
		name     string
		lat1     float64
		lon1     float64
		lat2     float64
		lon2     float64
		expected float64
	}{
		{
			name:     "Same point",
			lat1:     0.0,
			lon1:     0.0,
			lat2:     0.0,
			lon2:     0.0,
			expected: 0.0,
		},
		{
			name:     "Points on the equator",
			lat1:     0.0,
			lon1:     0.0,
			lat2:     0.0,
			lon2:     90.0,
			expected: 10007.54, // Approximate distance
		},
		{
			name:     "Points on the same meridian",
			lat1:     0.0,
			lon1:     0.0,
			lat2:     90.0,
			lon2:     0.0,
			expected: 10007.54, // Approximate distance
		},
		{
			name:     "Points with small distance",
			lat1:     52.2296756,
			lon1:     21.0122287,
			lat2:     52.406374,
			lon2:     16.9251681,
			expected: 278.546, // Approximate distance
		},
		{
			name:     "Points with large distance",
			lat1:     36.12,
			lon1:     -86.67,
			lat2:     33.94,
			lon2:     -118.40,
			expected: 2886.07, // Approximate distance
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HaversineDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2, 0, 0)
			if math.Abs(result-tt.expected) > 1.0 { // Allow a small margin of error
				t.Errorf("HaversineDistance() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIntersectsEdge(t *testing.T) {
	tests := []struct {
		name         string
		tileCenter   polygon.Point // Replacing edge with center point
		satellitePos polygon.Point
		tileRadiusKm float64
		altitude     float64
		expected     bool
	}{
		{
			name:         "Satellite within radius of tile center (start point)",
			tileCenter:   polygon.Point{Latitude: 0.0, Longitude: 0.0}, // Center of the tile is the start point
			satellitePos: polygon.Point{Latitude: 0.0, Longitude: 0.0},
			tileRadiusKm: 1.0,
			altitude:     0.0,
			expected:     true,
		},
		{
			name:         "Satellite within radius of tile center (end point)",
			tileCenter:   polygon.Point{Latitude: 0.5, Longitude: 0.5}, // Center of the tile is the midpoint between start and end
			satellitePos: polygon.Point{Latitude: 0.5, Longitude: 0.5},
			tileRadiusKm: 1.0,
			altitude:     0.0,
			expected:     true,
		},
		{
			name:         "Satellite outside radius of tile center",
			tileCenter:   polygon.Point{Latitude: 0.5, Longitude: 0.5},
			satellitePos: polygon.Point{Latitude: 2.0, Longitude: 2.0},
			tileRadiusKm: 1.0,
			altitude:     0.0,
			expected:     false,
		},
		{
			name:         "Satellite exactly 1 km away in latitude direction",
			tileCenter:   polygon.Point{Latitude: 0.0, Longitude: 0.0},
			satellitePos: polygon.Point{Latitude: 0.00899322, Longitude: 0.0}, // Exactly 1 km away along latitude
			tileRadiusKm: 1.0,
			altitude:     0.0,
			expected:     true,
		},
		{
			name:         "Satellite exactly 1 km away in longitude direction",
			tileCenter:   polygon.Point{Latitude: 0.0, Longitude: 0.0},
			satellitePos: polygon.Point{Latitude: 0.0, Longitude: 0.00899322}, // Exactly 1 km away along longitude
			tileRadiusKm: 1.0,
			altitude:     0.0,
			expected:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Calculate distance from satellite to the center of the tile
			centerDistance := HaversineDistance(tt.satellitePos.Latitude, tt.satellitePos.Longitude, tt.tileCenter.Latitude, tt.tileCenter.Longitude, 0, 0)

			// Log detailed information about the distance calculation
			t.Logf("[DEBUG] Satellite Position: (%.6f, %.6f)", tt.satellitePos.Latitude, tt.satellitePos.Longitude)
			t.Logf("[DEBUG] Tile Center: (%.6f, %.6f)", tt.tileCenter.Latitude, tt.tileCenter.Longitude)
			t.Logf("[DEBUG] Calculated Center Distance: %.6f km", centerDistance)
			t.Logf("[DEBUG] Tile Radius: %.6f km", tt.tileRadiusKm)

			// Adding margin of error for comparison
			marginOfError := 0.01
			t.Logf("[DEBUG] Margin of Error: %.6f km", marginOfError)

			// Check if satellite is within the tile radius from the center of the tile
			if got := Intersects(tt.tileCenter, tt.satellitePos, tt.tileRadiusKm, tt.altitude); got != tt.expected {
				t.Errorf("IntersectsEdge() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestComputeAOS(t *testing.T) {
	startTime := time.Now()
	endTime := startTime.Add(1 * time.Hour) // Extended time window
	timeStep := 5 * time.Second             // Reduced time step
	tileRadiusKm := 7000.0                  // Increased tile radius to 7000 km

	tleLine1 := "1 25544U 98067A   20344.54791435  .00001234  00000-0  29746-4 0  9998"
	tleLine2 := "2 25544  51.6456 212.9669 0001235 341.2074 106.3520 15.48921140255678"
	satrec := satellite.TLEToSat(tleLine1, tleLine2, satellite.GravityWGS84)

	tests := []struct {
		name          string
		vertices      []polygon.Point
		expectedAOS   bool
		expectedMaxEl bool
	}{
		{
			name: "Low Latitude Tile",
			vertices: []polygon.Point{
				{Latitude: 0.0, Longitude: 0.0},
				{Latitude: 0.0, Longitude: 10.0},
				{Latitude: 10.0, Longitude: 0.0},
				{Latitude: 10.0, Longitude: 10.0},
			},
			expectedAOS:   true,
			expectedMaxEl: true,
		},
		{
			name: "High Latitude Tile",
			vertices: []polygon.Point{
				{Latitude: 45.0, Longitude: 45.0},
				{Latitude: 45.0, Longitude: 55.0},
				{Latitude: 55.0, Longitude: 45.0},
				{Latitude: 55.0, Longitude: 55.0},
			},
			expectedAOS:   true,
			expectedMaxEl: true,
		},
		{
			name: "Out of Range Tile",
			vertices: []polygon.Point{
				{Latitude: 80.0, Longitude: 80.0},
				{Latitude: 80.0, Longitude: 90.0},
				{Latitude: 90.0, Longitude: 80.0},
				{Latitude: 90.0, Longitude: 90.0},
			},
			expectedAOS:   false,
			expectedMaxEl: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var maxElevation float64 = -1.0

			aos := ComputeAOS(satrec, tt.vertices[0], tileRadiusKm, startTime, endTime, timeStep, &maxElevation)

			if tt.expectedAOS && aos.IsZero() {
				t.Errorf("Expected AOS but got none.")
			} else if !tt.expectedAOS && !aos.IsZero() {
				t.Errorf("Did not expect AOS but got one.")
			}

			if tt.expectedMaxEl && maxElevation <= 0 {
				t.Errorf("Expected positive MaxElevation but got %.2f", maxElevation)
			} else if !tt.expectedMaxEl && maxElevation > 0 {
				t.Errorf("Did not expect MaxElevation but got %.2f", maxElevation)
			}
		})
	}
}

func TestComputeVisibilityWindow(t *testing.T) {
	startTime := time.Now()
	endTime := startTime.Add(24 * time.Hour) // Test over a 24-hour period
	timeStep := 1 * time.Minute              // 1-minute time step for granularity
	tileRadiusKm := 100.0                    // Test with a 100 km tile radius

	// Test case setup with TLE data
	tleLine1 := "1 25544U 98067A   20344.54791435  .00001234  00000-0  29746-4 0  9998"
	tleLine2 := "2 25544  51.6456 212.9669 0001235 341.2074 106.3520 15.48921140255678"

	// Vertices of the polygon (defining a square area for testing)
	vertices := []polygon.Point{
		{Latitude: 0.0, Longitude: 0.0},
		{Latitude: 0.0, Longitude: 10.0},
		{Latitude: 10.0, Longitude: 0.0},
		{Latitude: 10.0, Longitude: 10.0},
	}

	// Test cases setup
	tests := []struct {
		name                 string
		tleLine1             string
		tleLine2             string
		vertices             []polygon.Point
		tileRadiusKm         float64
		startTime            time.Time
		endTime              time.Time
		expectedAOS          time.Time
		expectedLOS          time.Time
		expectedMaxElevation float64
	}{
		{
			name:                 "Satellite visible within the tile radius",
			tleLine1:             tleLine1,
			tleLine2:             tleLine2,
			vertices:             vertices,
			tileRadiusKm:         tileRadiusKm,
			startTime:            startTime,
			endTime:              endTime,
			expectedAOS:          startTime.Add(5 * time.Minute),  // Assume AOS occurs after 5 minutes
			expectedLOS:          startTime.Add(10 * time.Minute), // Assume LOS occurs after 10 minutes
			expectedMaxElevation: 45.0,                            // Max elevation in degrees
		},
		{
			name:                 "Satellite outside tile radius",
			tleLine1:             tleLine1,
			tleLine2:             tleLine2,
			vertices:             vertices,
			tileRadiusKm:         tileRadiusKm,
			startTime:            startTime,
			endTime:              endTime,
			expectedAOS:          time.Time{}, // No AOS detected
			expectedLOS:          time.Time{}, // No LOS detected
			expectedMaxElevation: 0.0,
		},
		{
			name:                 "Satellite with low elevation",
			tleLine1:             tleLine1,
			tleLine2:             tleLine2,
			vertices:             vertices,
			tileRadiusKm:         tileRadiusKm,
			startTime:            startTime,
			endTime:              endTime,
			expectedAOS:          startTime.Add(10 * time.Minute), // Assume AOS occurs after 10 minutes
			expectedLOS:          startTime.Add(20 * time.Minute), // Assume LOS occurs after 20 minutes
			expectedMaxElevation: 5.0,                             // Low elevation
		},
		{
			name:                 "Satellite with high elevation",
			tleLine1:             tleLine1,
			tleLine2:             tleLine2,
			vertices:             vertices,
			tileRadiusKm:         tileRadiusKm,
			startTime:            startTime,
			endTime:              endTime,
			expectedAOS:          startTime.Add(15 * time.Minute), // Assume AOS occurs after 15 minutes
			expectedLOS:          startTime.Add(25 * time.Minute), // Assume LOS occurs after 25 minutes
			expectedMaxElevation: 85.0,                            // High elevation
		},
	}

	// Loop through all test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Compute the visibility window
			aos, los, maxElevation := ComputeVisibilityWindow("25544", tt.tleLine1, tt.tleLine2, tt.vertices[0], tt.tileRadiusKm, tt.startTime, tt.endTime, timeStep)

			// AOS Validation
			if !aos.Equal(tt.expectedAOS) {
				t.Errorf("Expected AOS at %v but got %v", tt.expectedAOS, aos)
			} else {
				t.Logf("AOS detected at %v", aos)
			}

			// LOS Validation
			if !los.Equal(tt.expectedLOS) {
				t.Errorf("Expected LOS at %v but got %v", tt.expectedLOS, los)
			} else {
				t.Logf("LOS detected at %v", los)
			}

			// Max Elevation Validation
			if maxElevation != tt.expectedMaxElevation {
				t.Errorf("Expected Max Elevation %.2f but got %.2f", tt.expectedMaxElevation, maxElevation)
			} else {
				t.Logf("Max Elevation detected: %.2f degrees", maxElevation)
			}
		})
	}
}

func TestCalculateIntegratedElevation(t *testing.T) {
	// Define test cases with satellite position, altitude, ground point, and expected elevation range
	testCases := []struct {
		name           string
		satellitePos   polygon.Point
		satelliteAltKm float64
		groundPoint    polygon.Point
	}{
		{
			name:           "Directly Overhead",
			satellitePos:   polygon.Point{Latitude: 0.0, Longitude: 0.0},
			satelliteAltKm: 500.0,
			groundPoint:    polygon.Point{Latitude: 0.0, Longitude: 0.0},
		},
		{
			name:           "Near Horizon",
			satellitePos:   polygon.Point{Latitude: 10.0, Longitude: 10.0},
			satelliteAltKm: 500.0,
			groundPoint:    polygon.Point{Latitude: 10.0, Longitude: 20.0},
		},
		{
			name:           "Far Distance",
			satellitePos:   polygon.Point{Latitude: 45.0, Longitude: 45.0},
			satelliteAltKm: 1000.0,
			groundPoint:    polygon.Point{Latitude: -45.0, Longitude: -45.0},
		},
	}

	// Iterate through test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			elevation := CalculateIntegratedElevationFromPoint(tc.satellitePos, tc.satelliteAltKm, tc.groundPoint)
			if elevation <= 0 {
				t.Errorf("For %s, expected positive elevation but got %.2f", tc.name, elevation)
			} else {
				t.Logf("For %s, calculated elevation: %.2f degrees", tc.name, elevation)
			}
		})
	}
}

func TestGeneratePolygonEdges(t *testing.T) {
	tests := []struct {
		name     string
		vertices []polygon.Point
		expected []polygon.Edge
	}{
		{
			name: "Square Polygon",
			vertices: []polygon.Point{
				{Latitude: 0.0, Longitude: 0.0},
				{Latitude: 0.0, Longitude: 1.0},
				{Latitude: 1.0, Longitude: 1.0},
				{Latitude: 1.0, Longitude: 0.0},
			},
			expected: []polygon.Edge{
				{Start: polygon.Point{Latitude: 0.0, Longitude: 0.0}, End: polygon.Point{Latitude: 0.0, Longitude: 1.0}},
				{Start: polygon.Point{Latitude: 0.0, Longitude: 1.0}, End: polygon.Point{Latitude: 1.0, Longitude: 1.0}},
				{Start: polygon.Point{Latitude: 1.0, Longitude: 1.0}, End: polygon.Point{Latitude: 1.0, Longitude: 0.0}},
				{Start: polygon.Point{Latitude: 1.0, Longitude: 0.0}, End: polygon.Point{Latitude: 0.0, Longitude: 0.0}},
			},
		},
		{
			name: "Triangle Polygon",
			vertices: []polygon.Point{
				{Latitude: 0.0, Longitude: 0.0},
				{Latitude: 1.0, Longitude: 1.0},
				{Latitude: 0.0, Longitude: 2.0},
			},
			expected: []polygon.Edge{
				{Start: polygon.Point{Latitude: 0.0, Longitude: 0.0}, End: polygon.Point{Latitude: 1.0, Longitude: 1.0}},
				{Start: polygon.Point{Latitude: 1.0, Longitude: 1.0}, End: polygon.Point{Latitude: 0.0, Longitude: 2.0}},
				{Start: polygon.Point{Latitude: 0.0, Longitude: 2.0}, End: polygon.Point{Latitude: 0.0, Longitude: 0.0}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			edges := GeneratePolygonEdges(tt.vertices)
			if len(edges) != len(tt.expected) {
				t.Fatalf("Expected %d edges but got %d", len(tt.expected), len(edges))
			}
			for i, edge := range edges {
				if edge != tt.expected[i] {
					t.Errorf("Edge mismatch at index %d: expected %v, got %v", i, tt.expected[i], edge)
				}
			}
		})
	}
}
