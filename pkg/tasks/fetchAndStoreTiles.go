package tasks

import (
	"fmt"
	"log"
	"math"

	"github.com/Elbujito/2112/pkg/api/handlers/tiles"
	"github.com/Elbujito/2112/pkg/db/models"
	"gorm.io/gorm"
)

func init() {
	// Initialize the fetch and store tiles task
	task := &Task{
		Name:         "fetchAndStoreTiles",
		Description:  "Fetches tiles and stores them in the database",
		RequiredArgs: []string{}, // No required arguments
		Run:          execFetchAndStoreTilesTask,
	}
	Tasks.AddTask(task)
}

// Task that will fetch and store tiles
func execFetchAndStoreTilesTask(env *TaskEnv, args map[string]string) error {
	// Define the zoom level and geographical bounds for the area you want to cover
	zoomLevel := 10                  // Define the zoom level you want to use
	latMin, latMax := 37.0, 38.0     // Example latitude range (e.g., San Francisco to Oakland)
	lonMin, lonMax := -123.0, -122.0 // Example longitude range

	// Calculate the number of tiles required to cover this area at the given zoom level
	tileWidth, tileHeight, totalTiles := calculateTileDimensions(latMin, latMax, lonMin, lonMax, zoomLevel)

	// Limit to generating 1000 tiles
	if totalTiles > 1000 {
		// Adjust the range to ensure no more than 1000 tiles are generated
		scaleFactor := math.Sqrt(float64(1000) / float64(totalTiles))
		latRange := (latMax - latMin) * scaleFactor
		lonRange := (lonMax - lonMin) * scaleFactor

		// Recalculate the number of tiles for the adjusted area
		tileWidth, tileHeight, _ = calculateTileDimensions(latMin, latMin+latRange, lonMin, lonMin+lonRange, zoomLevel)
	}

	log.Printf("Using zoom level %d, generating %d tiles", zoomLevel, tileWidth*tileHeight)

	// Get the TileService from the models package
	tileService := models.TileModel()

	// Iterate over the coordinate range and fetch/store tiles
	for x := 0; x < tileWidth; x++ {
		for y := 0; y < tileHeight; y++ {
			// Convert lat, lon to tile x, y and fetch the tile data
			lat, lon := calculateTileCenter(latMin, lonMin, x, y, zoomLevel)
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

// calculateTileDimensions calculates the number of tiles needed to cover a specified geographical area
func calculateTileDimensions(latMin, latMax, lonMin, lonMax float64, zoom int) (int, int, int) {
	// Calculate the number of tiles in both x and y directions for a given zoom level
	tileWidth := int(math.Pow(2, float64(zoom)) * (lonMax - lonMin) / 360)
	tileHeight := int(math.Pow(2, float64(zoom)) * (latMax - latMin) / 180)

	// Total number of tiles is the product of the width and height in tiles
	totalTiles := tileWidth * tileHeight
	return tileWidth, tileHeight, totalTiles
}

// calculateTileCenter calculates the center latitude and longitude for a specific tile
func calculateTileCenter(latMin, lonMin float64, x int, y int, zoom int) (float64, float64) {
	// Number of tiles at the zoom level
	n := math.Pow(2, float64(zoom))

	// Calculate the longitude of the top-left corner of the tile
	lon := float64(x)/n*360.0 - 180.0

	// Calculate the latitude of the top-left corner of the tile using the Mercator projection
	latRad := math.Atan(math.Sinh(math.Pi * (1 - 2*float64(y)/n)))
	lat := latRad * 180.0 / math.Pi

	// Now we calculate the center latitude and longitude for the tile
	// The range of latitudes and longitudes are scaled by x and y coordinates
	centerLat := latMin + (lat - latMin) // latitude center
	centerLon := lonMin + (lon - lonMin) // longitude center

	return centerLat, centerLon
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
