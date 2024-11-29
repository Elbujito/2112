package handlers

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/Elbujito/2112/internal/domain"
	"github.com/Elbujito/2112/pkg/fx/polygon"
)

type TileProvisionHandler struct {
	tileRepo domain.TileRepository
}

// NewTileProvisionHandler creates a new instance of TileProvisionHandler.
func NewTileProvisionHandler(tileRepo domain.TileRepository) TileProvisionHandler {
	return TileProvisionHandler{
		tileRepo: tileRepo,
	}
}

// GetTask provides metadata about this handler's task.
func (h *TileProvisionHandler) GetTask() Task {
	return Task{
		Name:        "fetchAndStoreTiles",
		Description: "Fetches tiles and stores them in the database",
		RequiredArgs: []string{
			"radiusInMeter",
			"faces",
		},
	}
}

// Run executes the handler's task with the provided arguments.
func (h *TileProvisionHandler) Run(ctx context.Context, args map[string]string) error {
	// Parse arguments
	radiusInMeter, ok := args["radiusInMeter"]
	if !ok || radiusInMeter == "" {
		return fmt.Errorf("missing required argument: radiusInMeter")
	}

	nbFaces, ok := args["faces"]
	if !ok || nbFaces == "" {
		return fmt.Errorf("missing required argument: faces")
	}

	radius, err := strconv.ParseFloat(radiusInMeter, 64)
	if err != nil {
		return fmt.Errorf("invalid radiusInMeter: %w", err)
	}

	faces, err := strconv.Atoi(nbFaces)
	if err != nil {
		return fmt.Errorf("invalid faces value: %w", err)
	}

	// Generate all tiles for the given radius and number of faces
	polygons := polygon.GenerateAllTilesForRadius(radius, faces)

	// Batch size for database operations
	const batchSize = 100
	tileChan := make(chan domain.Tile, len(polygons))
	errChan := make(chan error, len(polygons))
	var wg sync.WaitGroup

	// Concurrent workers for batch processing
	numWorkers := 4
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for tile := range tileChan {
				// Batch Upsert operation
				if err := h.tileRepo.Upsert(ctx, tile); err != nil {
					errChan <- fmt.Errorf("failed to upsert tile with Quadkey %s: %w", tile.Quadkey, err)
				}
			}
		}()
	}

	// Feed tiles into the worker channel
	go func() {
		for _, p := range polygons {
			tile := domain.NewTile(p)
			tileChan <- tile
		}
		close(tileChan)
	}()

	// Wait for workers to finish
	wg.Wait()
	close(errChan)

	// Collect errors
	var combinedErr error
	for err := range errChan {
		combinedErr = fmt.Errorf("%v; %w", combinedErr, err)
	}

	return combinedErr
}
