package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Elbujito/2112/internal/domain"
	"github.com/Elbujito/2112/pkg/fx/polygon"
)

type TileProvisionHandler struct {
	tileRepo domain.TileRepository
}

func NewTileProvisionHandler(
	tileRepo domain.TileRepository) TileProvisionHandler {
	return TileProvisionHandler{
		tileRepo: tileRepo,
	}
}

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

func (h *TileProvisionHandler) Run(ctx context.Context, args map[string]string) error {
	radiusInMeter, ok := args["radiusInMeter"]
	if !ok || radiusInMeter == "" {
		return fmt.Errorf("missing required argument: radius")
	}

	nbFaces, ok := args["faces"]
	if !ok || nbFaces == "" {
		return fmt.Errorf("missing required argument: nbFaces")
	}

	radius, err := strconv.ParseFloat(radiusInMeter, 64)
	if err != nil {
		return err
	}

	faces, err := strconv.Atoi(nbFaces)
	if err != nil {
		return err
	}

	polygons := polygon.GenerateAllTilesForRadius(radius, faces)
	for _, p := range polygons {
		tile := domain.NewTile(p)
		err := h.tileRepo.Upsert(ctx, tile)
		if err != nil {
			return fmt.Errorf("failed to upsert Tile for Key %s: %v", tile.Quadkey, err)
		}
	}
	return nil
}
