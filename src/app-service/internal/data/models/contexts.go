package models

import (
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/domain"
)

// Context represents the database model for logical groupings.
type Context struct {
	ModelBase
	Name        string `gorm:"size:255;unique;not null"` // Unique name of the context
	Description string `gorm:"size:1024"`                // Optional description of the context
	ActivatedAt *time.Time
}

// MapToContextDomain converts a Context database model to a Context domain model.
func MapToContextDomain(c Context) domain.Context {
	return domain.Context{
		ModelBase: domain.ModelBase{
			ID:          c.ID,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   &c.UpdatedAt,
			DeleteAt:    c.DeleteAt,
			ProcessedAt: c.ProcessedAt,
			IsActive:    c.IsActive,
			IsFavourite: c.IsFavourite,
			DisplayName: c.DisplayName,
		},
		Name:        c.Name,
		Description: c.Description,
		ActivatedAt: c.ActivatedAt,
	}
}

// MapToContextModel converts a Context domain model to a Context database model.
func MapToContextModel(c domain.Context) Context {
	return Context{
		ModelBase: ModelBase{
			ID:          c.ModelBase.ID,
			CreatedAt:   c.ModelBase.CreatedAt,
			UpdatedAt:   *c.ModelBase.UpdatedAt,
			DeleteAt:    c.ModelBase.DeleteAt,
			ProcessedAt: c.ModelBase.ProcessedAt,
			IsActive:    c.ModelBase.IsActive,
			IsFavourite: c.ModelBase.IsFavourite,
			DisplayName: c.ModelBase.DisplayName,
		},
		Name:        c.Name,
		Description: c.Description,
		ActivatedAt: c.ActivatedAt,
	}
}

// ContextTile defines the many-to-many relationship between Context and Tile.
type ContextTile struct {
	ContextID string  `gorm:"not null;index"` // Foreign key to Context
	TileID    string  `gorm:"not null;index"` // Foreign key to Tile
	Context   Context `gorm:"constraint:OnDelete:CASCADE;foreignKey:ContextID;references:ID"`
	Tile      Tile    `gorm:"constraint:OnDelete:CASCADE;foreignKey:TileID;references:ID"`
}

// MapToContextTileDomain converts a ContextTile database model to a ContextTile domain model.
func MapToContextTileDomain(ct ContextTile) domain.ContextTile {
	return domain.ContextTile{
		ContextID: ct.ContextID,
		TileID:    ct.TileID,
	}
}

// MapToContextTileModel converts a ContextTile domain model to a ContextTile database model.
func MapToContextTileModel(ct domain.ContextTile) ContextTile {
	return ContextTile{
		ContextID: ct.ContextID,
		TileID:    ct.TileID,
	}
}

// TableName overrides the default table name for ContextTile.
func (ContextTile) TableName() string {
	return "context_tiles"
}

// ContextTLE defines the many-to-many relationship between Context and TLE.
type ContextTLE struct {
	ContextID string  `gorm:"not null;index"` // Foreign key to Context
	TLEID     string  `gorm:"not null;index"` // Foreign key to TLE
	Context   Context `gorm:"constraint:OnDelete:CASCADE;foreignKey:ContextID;references:ID"`
	TLE       TLE     `gorm:"constraint:OnDelete:CASCADE;foreignKey:TLEID;references:ID"`
}

// MapToContextTLEDomain converts a ContextTLE database model to a ContextTLE domain model.
func MapToContextTLEDomain(ct ContextTLE) domain.ContextTLE {
	return domain.ContextTLE{
		ContextID: ct.ContextID,
		TLEID:     ct.TLEID,
	}
}

// MapToContextTLEModel converts a ContextTLE domain model to a ContextTLE database model.
func MapToContextTLEModel(ct domain.ContextTLE) ContextTLE {
	return ContextTLE{
		ContextID: ct.ContextID,
		TLEID:     ct.TLEID,
	}
}

// TableName overrides the default table name for ContextTLE.
func (ContextTLE) TableName() string {
	return "context_tles"
}

// ContextSatellite defines the many-to-many relationship between Context and Satellite.
type ContextSatellite struct {
	ContextID   string    `gorm:"not null;index"` // Foreign key to Context
	SatelliteID string    `gorm:"not null;index"` // Foreign key to Satellite
	Context     Context   `gorm:"constraint:OnDelete:CASCADE;foreignKey:ContextID;references:ID"`
	Satellite   Satellite `gorm:"constraint:OnDelete:CASCADE;foreignKey:SatelliteID;references:ID"`
}

// MapToContextSatelliteDomain converts a ContextSatellite database model to a ContextSatellite domain model.
func MapToContextSatelliteDomain(cs ContextSatellite) domain.ContextSatellite {
	return domain.ContextSatellite{
		ContextID:   cs.ContextID,
		SatelliteID: cs.SatelliteID,
	}
}

// MapToContextSatelliteModel converts a ContextSatellite domain model to a ContextSatellite database model.
func MapToContextSatelliteModel(cs domain.ContextSatellite) ContextSatellite {
	return ContextSatellite{
		ContextID:   cs.ContextID,
		SatelliteID: cs.SatelliteID,
	}
}

// TableName overrides the default table name for ContextSatellite.
func (ContextSatellite) TableName() string {
	return "context_satellites"
}
