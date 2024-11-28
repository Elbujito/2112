package polygon

import (
	"log"
	"math"
	"sync"

	"github.com/Elbujito/2112/pkg/fx/constants"
)

// Polygon represents a geographic polygon
type Polygon struct {
	Center     Quadkey
	NbFaces    int
	Radius     float64
	Boundaries []Quadkey
}

// NewPolygon creates a new Polygon given the number of faces, center coordinates, zoom level, and radius
func NewPolygon(nbFaces int, center LatLong, level int, radius float64) Polygon {
	boundaries := generateBoundaries(nbFaces, center, level, radius)
	return Polygon{
		NbFaces:    nbFaces,
		Radius:     radius,
		Center:     NewQuadkey(center.LatDegrees(), center.LonDegrees(), level),
		Boundaries: boundaries,
	}
}

// generateBoundaries calculates the boundary points of the polygon
func generateBoundaries(nbFaces int, center LatLong, level int, radius float64) []Quadkey {
	boundaries := make([]Quadkey, 0, nbFaces)

	centerLatRad := center.LatRadians()
	centerLonRad := center.LonRadians()

	angleStep := 2 * math.Pi / float64(nbFaces)

	for i := 0; i < nbFaces; i++ {
		angle := angleStep * float64(i)

		deltaLat := radius * math.Cos(angle) / constants.EARTH_RADIUS
		deltaLon := radius * math.Sin(angle) / (constants.EARTH_RADIUS * math.Cos(centerLatRad))

		newLat := centerLatRad + deltaLat
		newLon := centerLonRad + deltaLon

		newLatDeg := newLat * 180 / math.Pi
		newLonDeg := newLon * 180 / math.Pi

		boundary := NewQuadkey(newLatDeg, newLonDeg, level)
		boundaries = append(boundaries, boundary)
	}

	return boundaries
}

func calculateZoomLevelForTileRadius(tileRadius float64) int {
	const earthCircumference = constants.EARTH_CIRCUMFERENCE_METER
	for zoom := 0; zoom <= constants.MAX_ZOOM_LEVEL; zoom++ {
		tileSize := earthCircumference / math.Pow(2, float64(zoom))

		log.Printf("Zoom: %d, Tile Size: %.2f, Half Tile Size: %.2f", zoom, tileSize, tileSize/2)
		if tileSize/2 <= tileRadius {
			return zoom
		}
	}
	return constants.MAX_ZOOM_LEVEL
}

// GenerateAllTilesForRadius generates all tiles and their polygons for a given radius
func GenerateAllTilesForRadius(tileRadius float64, nbFaces int) map[Quadkey]Polygon {
	zoom := calculateZoomLevelForTileRadius(tileRadius)
	numTiles := int(math.Pow(2, float64(zoom)))

	tilePolygons := make(map[Quadkey]Polygon)
	tileChan := make(chan struct {
		Quadkey Quadkey
		Polygon Polygon
	})

	var wg sync.WaitGroup
	worker := func(startX, endX int) {
		defer wg.Done()
		for x := startX; x < endX; x++ {
			for y := 0; y < numTiles; y++ {
				lat, lon := TileXYToLatLon(x, y, zoom)
				center := LatLong{Lat: Coordinate{lat}, Lon: Coordinate{lon}}
				centerQuadkey := NewQuadkey(center.LatDegrees(), center.LonDegrees(), zoom)
				polygon := NewPolygon(nbFaces, center, zoom, tileRadius)

				tileChan <- struct {
					Quadkey Quadkey
					Polygon Polygon
				}{centerQuadkey, polygon}
			}
		}
	}

	numWorkers := 4
	batchSize := numTiles / numWorkers
	for i := 0; i < numWorkers; i++ {
		startX := i * batchSize
		endX := startX + batchSize
		if i == numWorkers-1 {
			endX = numTiles
		}

		wg.Add(1)
		go worker(startX, endX)
	}

	go func() {
		wg.Wait()
		close(tileChan)
	}()

	for result := range tileChan {
		tilePolygons[result.Quadkey] = result.Polygon
	}

	return tilePolygons
}

// TileXYToLatLon converts tile coordinates to lat/lon
func TileXYToLatLon(x, y, zoom int) (float64, float64) {
	n := math.Pow(2, float64(zoom))
	lon := float64(x)/n*360.0 - 180.0

	latRad := math.Atan(math.Sinh(math.Pi * (1 - 2*float64(y)/n)))
	lat := latRad * 180.0 / math.Pi

	return lat, lon
}

// LatLonToTileXY converts lat/lon to tile coordinates
func LatLonToTileXY(lat, lon float64, zoom int) (int, int) {
	n := math.Pow(2, float64(zoom))

	tileX := int((lon + 180.0) / 360.0 * n)
	tileY := int((1.0 - math.Log(math.Tan(lat*math.Pi/180.0)+1.0/math.Cos(lat*math.Pi/180.0))/math.Pi) / 2.0 * n)

	return tileX, tileY
}
