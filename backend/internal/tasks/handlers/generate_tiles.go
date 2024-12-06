package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Elbujito/2112/internal/domain"
)

type GenerateTilesHandler struct {
	tileRepo domain.TileRepository
}

// NewGenerateTilesHandler creates a new instance of TileProvisionHandler.
func NewGenerateTilesHandler(tileRepo domain.TileRepository) GenerateTilesHandler {
	return GenerateTilesHandler{
		tileRepo: tileRepo,
	}
}

// GetTask provides metadata about this handler's task.
func (h *GenerateTilesHandler) GetTask() Task {
	return Task{
		Name:        "generate_tiles",
		Description: "Fetches tiles and stores them in the database",
		RequiredArgs: []string{
			"radiusInMeter",
			"faces",
		},
	}
}

// Run executes the handler's task with the provided arguments.
func (h *GenerateTilesHandler) Run(ctx context.Context, args map[string]string) error {
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

	polygons := xpolygon.GenerateAllTilesForRadius(radius, faces)

	for _, p := range polygons {
		tile := domain.NewTile(p)
		err := h.tileRepo.Upsert(ctx, tile)
		if err != nil {
			return fmt.Errorf("failed to upsert Tile for Key %s: %v", tile.Quadkey, err)
		}
	}
	return nil
}
