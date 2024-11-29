package migrations

import (
	"time"

	"github.com/Elbujito/2112/internal/data/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "2024112604_create_quadkey_visibility_schema",
		Migrate: func(db *gorm.DB) error {
			type Satellite struct {
				models.ModelBase
				Name    string `gorm:"size:255;not null"`
				NoradID string `gorm:"size:255;unique;not null"`
				Type    string `gorm:"size:255"` // Optional satellite type
			}

			type TLE struct {
				models.ModelBase
				NoradID string    `gorm:"not null;index"` // Foreign key to Satellite table
				Line1   string    `gorm:"size:255;not null"`
				Line2   string    `gorm:"size:255;not null"`
				Epoch   time.Time `gorm:"not null"` // Time associated with the TLE
			}

			type Tile struct {
				models.ModelBase
				Quadkey        string  `gorm:"size:256;unique;not null"` // Unique identifier for the tile (Quadkey)
				ZoomLevel      int     `gorm:"not null"`                 // Zoom level for the tile
				CenterLat      float64 `gorm:"not null"`                 // Center latitude of the tile
				CenterLon      float64 `gorm:"not null"`                 // Center longitude of the tile
				NbFaces        int     `gorm:"not null"`                 // Number of faces in the tile's shape
				Radius         float64 `gorm:"not null"`                 // Radius of the tile in meters
				BoundariesJSON string  `gorm:"type:json"`                // Serialized JSON of the boundary vertices of the tile
				SpatialIndex   string  `gorm:"type:geometry;spatialIndex"`
			}

			type TileSatelliteMapping struct {
				models.ModelBase
				NoradID      string    `gorm:"not null;index"` // Foreign key to Satellite table
				TileID       string    `gorm:"not null;index"` // Foreign key to Tile table
				Tile         Tile      `gorm:"constraint:OnDelete:CASCADE;foreignKey:TileID;references:ID"`
				Aos          time.Time `gorm:"not null"` // Visibility start time
				MaxElevation float64   `gorm:"not null"` // Max elevation during visibility in degrees
			}

			// AutoMigrate all tables
			return db.AutoMigrate(&Satellite{}, &TLE{}, &Tile{}, &TileSatelliteMapping{})
		},
		Rollback: func(db *gorm.DB) error {
			return db.Migrator().DropTable("visibilities", "tiles", "tles", "satellites")
		},
	}

	AddMigration(m)
}
