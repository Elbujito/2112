package tasks

import (
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/Elbujito/2112/pkg/api/handlers/tiles"
	"github.com/Elbujito/2112/pkg/db/models"
	"gorm.io/gorm"
)

func init() {
	// Initialize the fetch and store tiles task
	task := &Task{
		Name:        "fetchAndStoreTiles",
		Description: "Fetches tiles and stores them in the database",
		RequiredArgs: []string{
			"desiredTiles", // Number of tiles to generate
		},
		Run: execFetchAndStoreTilesTask,
	}
	Tasks.AddTask(task)
}

func execFetchAndStoreTilesTask(env *TaskEnv, args map[string]string) error {
	// Parse the desiredTiles argument
	desiredTilesStr, ok := args["desiredTiles"]
	if !ok {
		return fmt.Errorf("required argument 'desiredTiles' is missing")
	}

	desiredTiles, err := strconv.Atoi(desiredTilesStr)
	if err != nil {
		return fmt.Errorf("invalid 'desiredTiles' value: %v", err)
	}

	// Define the initial geographical bounds (entire globe)
	latMin, latMax := -85.0511, 85.0511
	lonMin, lonMax := -180.0, 180.0

	// Dynamically calculate the appropriate zoom level to meet the desired number of tiles
	zoomLevel := calculateZoomLevelForGlobalCoverage(desiredTiles)

	// Calculate the number of tiles at the calculated zoom level
	tileWidth, tileHeight, totalTiles := calculateTileDimensions(latMin, latMax, lonMin, lonMax, zoomLevel)

	log.Printf("Using zoom level %d, generating %d tiles (%dx%d)", zoomLevel, totalTiles, tileWidth, tileHeight)

	// Get the TileService from the models package
	tileService := models.TileModel()

	// Iterate over the tile grid and fetch/store tiles
	for x := 0; x < tileWidth; x++ {
		for y := 0; y < tileHeight; y++ {
			// Convert tile coordinates to lat/lon and fetch the tile data
			lat, lon := tileXYToLatLon(x, y, zoomLevel)
			err := fetchTileAndStore(zoomLevel, lat, lon, tileService)
			if err != nil {
				log.Printf("Failed to fetch and store tile at zoom %d, lat %f, lon %f: %v", zoomLevel, lat, lon, err)
			} else {
				log.Printf("Successfully fetched and stored tile at zoom %d, lat %f, lon %f", zoomLevel, lat, lon)
			}
		}
	}

	return nil
}

// calculateZoomLevelForGlobalCoverage calculates the zoom level needed for the desired number of tiles
func calculateZoomLevelForGlobalCoverage(desiredTiles int) int {
	for zoom := 0; zoom <= 21; zoom++ {
		// Total tiles at zoom level = (2^zoom) * (2^zoom)
		totalTiles := int(math.Pow(2, float64(zoom)) * math.Pow(2, float64(zoom)))
		if totalTiles >= desiredTiles {
			return zoom
		}
	}
	return 21 // Default to max zoom if desiredTiles is very high
}

func calculateTileDimensions(latMin, latMax, lonMin, lonMax float64, zoom int) (int, int, int) {
	tileXMin, tileYMin := latLonToTileXY(latMin, lonMin, zoom)
	tileXMax, tileYMax := latLonToTileXY(latMax, lonMax, zoom)

	// Ensure the tile ranges are ordered correctly
	if tileXMax < tileXMin {
		tileXMin, tileXMax = tileXMax, tileXMin
	}
	if tileYMax < tileYMin {
		tileYMin, tileYMax = tileYMax, tileYMin
	}

	tileWidth := tileXMax - tileXMin + 1
	tileHeight := tileYMax - tileYMin + 1
	totalTiles := tileWidth * tileHeight

	return tileWidth, tileHeight, totalTiles
}

// tileXYToLatLon converts tile coordinates to lat/lon
func tileXYToLatLon(x, y, zoom int) (float64, float64) {
	n := math.Pow(2, float64(zoom))
	lon := float64(x)/n*360.0 - 180.0

	latRad := math.Atan(math.Sinh(math.Pi * (1 - 2*float64(y)/n)))
	lat := latRad * 180.0 / math.Pi

	return lat, lon
}

// latLonToTileXY converts lat/lon to tile coordinates
func latLonToTileXY(lat, lon float64, zoom int) (int, int) {
	n := math.Pow(2, float64(zoom))

	tileX := int((lon + 180.0) / 360.0 * n)
	tileY := int((1.0 - math.Log(math.Tan(lat*math.Pi/180.0)+1.0/math.Cos(lat*math.Pi/180.0))/math.Pi) / 2.0 * n)

	return tileX, tileY
}

// calculateTileCenter calculates the center latitude and longitude for a specific tile
func calculateTileCenter(latMax, lonMax, latMin, lonMin float64, x, y, tileWidth, tileHeight int) (float64, float64) {
	// Calculate the size of each tile in degrees for latitude and longitude
	tileLatSize := (latMax - latMin) / float64(tileHeight)
	tileLonSize := (lonMax - lonMin) / float64(tileWidth)

	// Compute the center of the tile at (x, y)
	lat := latMin + float64(y)*tileLatSize + tileLatSize/2
	lon := lonMin + float64(x)*tileLonSize + tileLonSize/2

	return lat, lon
}

// fetchTileAndStore fetches the tile and stores it in the database
func fetchTileAndStore(zoom int, lat float64, lon float64, tileService models.TileService) error {
	// Fetch the tile from the server
	tile, err := tiles.FetchTile(zoom, lat, lon) // Using the FetchTile function defined in the handler package
	if err != nil {
		return fmt.Errorf("failed to fetch tile: %v", err)
	}

	// Generate a unique Quadkey for this tile (you can adjust this based on your needs)
	quadkey := generateQuadkey(zoom, lat, lon)

	// Store the tile in the database using TileService
	return upsertTileInDB(quadkey, zoom, lat, lon, tile.TileData, tileService)
}

// Generate a Quadkey from zoom, lat, and lon
func generateQuadkey(zoom int, lat float64, lon float64) string {
	// Generate a quadkey from zoom, lat, and lon (this is just a simple example)
	return fmt.Sprintf("%d-%f-%f", zoom, lat, lon)
}

// upsertTileInDB checks if a tile exists and updates it, or inserts it if it doesn't exist
func upsertTileInDB(quadkey string, zoom int, lat float64, lon float64, tileData []byte, tileService models.TileService) error {
	// Check if the tile already exists using TileService's FindByQuadkey
	existingTile, err := tileService.FindByQuadkey(quadkey)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to check if tile exists: %v", err)
	}

	// If the tile exists, update it
	if existingTile != nil {
		existingTile.Quadkey = quadkey
		existingTile.ZoomLevel = zoom
		existingTile.CenterLat = lat
		existingTile.CenterLon = lon
		return tileService.Update(existingTile) // Use the Update method to update the existing tile
	}

	// If the tile doesn't exist, insert a new one
	newTile := models.Tile{
		Quadkey:   quadkey,
		ZoomLevel: zoom,
		CenterLat: lat,
		CenterLon: lon,
	}
	return newTile.Create() // Use the Create method to insert a new tile
}
