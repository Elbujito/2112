package polygon

import (
	"log"
	"testing"
	"time"
)

// Constants for test cases
const (
	TestTileRadius = 10000.0
	TestNbFaces    = 6
	TestZoomLevel  = 15
)

// TestNewPolygon verifies the creation of a new Polygon
func TestNewPolygon(t *testing.T) {
	center := LatLong{
		Lat: Coordinate{c: 40.7128},
		Lon: Coordinate{c: -74.0060},
	}

	polygon := NewPolygon(TestNbFaces, center, TestZoomLevel, TestTileRadius)

	if polygon.NbFaces != TestNbFaces {
		t.Errorf("Expected nbFaces %d, got %d", TestNbFaces, polygon.NbFaces)
	}
	if polygon.Radius != TestTileRadius {
		t.Errorf("Expected radius %.2f, got %.2f", TestTileRadius, polygon.Radius)
	}
	if len(polygon.Boundaries) != TestNbFaces {
		t.Errorf("Expected %d boundaries, got %d", TestNbFaces, len(polygon.Boundaries))
	}
}

// TestGenerateBoundaries verifies that boundaries are calculated correctly
func TestGenerateBoundaries(t *testing.T) {
	center := LatLong{
		Lat: Coordinate{c: 0.0},
		Lon: Coordinate{c: 0.0},
	}

	boundaries := generateBoundaries(TestNbFaces, center, TestZoomLevel, TestTileRadius)

	if len(boundaries) != TestNbFaces {
		t.Errorf("Expected %d boundaries, got %d", TestNbFaces, len(boundaries))
	}
}

// TestCalculateZoomLevelForTileRadius tests the zoom level calculation
func TestCalculateZoomLevelForTileRadius(t *testing.T) {
	tileRadius := 500.0

	expectedZoom := 16
	calculatedZoom := calculateZoomLevelForTileRadius(tileRadius)

	if calculatedZoom != expectedZoom {
		t.Errorf("Expected zoom level %d, got %d", expectedZoom, calculatedZoom)
	}
}

func TestTileXYToLatLon(t *testing.T) {
	zoom := 15
	x, y := 0, 0

	lat, lon := TileXYToLatLon(x, y, zoom)
	if lat > 85.05112878 || lat < -85.05112878 {
		t.Errorf("Latitude out of bounds: %.6f", lat)
	}
	if lon > 180 || lon < -180 {
		t.Errorf("Longitude out of bounds: %.6f", lon)
	}
}

// TestLatLonToTileXY tests the LatLonToTileXY function
func TestLatLonToTileXY(t *testing.T) {
	lat, lon := 40.7128, -74.0060
	zoom := 15

	x, y := LatLonToTileXY(lat, lon, zoom)
	if x < 0 || y < 0 {
		t.Errorf("Invalid tile coordinates: x=%d, y=%d", x, y)
	}
}

// TestGenerateAllTilesForRadius tests the GenerateAllTilesForRadius function
func TestGenerateAllTilesForRadius(t *testing.T) {
	start := time.Now()
	tiles := GenerateAllTilesForRadius(TestTileRadius, TestNbFaces)
	duration := time.Since(start)
	log.Printf("Generated %d tiles in %s", len(tiles), duration)
	if len(tiles) == 0 {
		t.Error("Expected non-empty map of tiles, got empty map")
	}
	for quadkey, polygon := range tiles {
		if polygon.NbFaces != TestNbFaces {
			t.Errorf("Expected polygon with %d faces, got %d for quadkey %v", TestNbFaces, polygon.NbFaces, quadkey)
		}
		break
	}
}
