package migrations

import (
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/data/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "2025011102_init",
		Migrate: func(db *gorm.DB) error {
			// Define the Context table
			type Context struct {
				models.ModelBase
				Name                       string     `gorm:"size:255;unique;not null"` // Context name
				TenantID                   string     `gorm:"size:255;not null;index"`  // Tenant identifier
				Description                string     `gorm:"size:1024"`                // Optional description
				MaxSatellite               int        `gorm:"not null"`
				MaxTiles                   int        `gorm:"not null"`
				ActivatedAt                *time.Time `gorm:"null"`
				DesactivatedAt             *time.Time `gorm:"null"`
				TriggerGeneratedMappingAt  *time.Time `gorm:"null"`
				TriggerImportedTLEAt       *time.Time `gorm:"null"`
				TriggerImportedSatelliteAt *time.Time `gorm:"null"`
			}

			// Define the Satellite table
			type Satellite struct {
				models.ModelBase
				Name           string     `gorm:"size:255;not null"`
				NoradID        string     `gorm:"size:255;unique;not null"`
				Type           string     `gorm:"size:255"`
				LaunchDate     *time.Time `gorm:"type:date"`
				DecayDate      *time.Time `gorm:"type:date"`
				IntlDesignator string     `gorm:"size:255"`
				Owner          string     `gorm:"size:255"`
				ObjectType     string     `gorm:"size:255"`
				Period         *float64   `gorm:"type:float"`
				Inclination    *float64   `gorm:"type:float"`
				Apogee         *float64   `gorm:"type:float"`
				Perigee        *float64   `gorm:"type:float"`
				RCS            *float64   `gorm:"type:float"`
				Altitude       *float64   `gorm:"type:float"`
			}

			// Define the TLE table
			type TLE struct {
				models.ModelBase
				NoradID string    `gorm:"not null;index"`
				Line1   string    `gorm:"size:255;not null"`
				Line2   string    `gorm:"size:255;not null"`
				Epoch   time.Time `gorm:"not null"`
			}

			// Define the Tile table
			type Tile struct {
				models.ModelBase
				Quadkey        string  `gorm:"size:256;unique;not null"`
				ZoomLevel      int     `gorm:"not null"`
				CenterLat      float64 `gorm:"not null"`
				CenterLon      float64 `gorm:"not null"`
				NbFaces        int     `gorm:"not null"`
				Radius         float64 `gorm:"not null"`
				BoundariesJSON string  `gorm:"type:json"`
				SpatialIndex   string  `gorm:"type:geometry(Polygon, 4326);spatialIndex"`
			}

			// Define the TileSatelliteMapping table
			type TileSatelliteMapping struct {
				models.ModelBase
				NoradID               string    `gorm:"not null;index"`
				TileID                string    `gorm:"not null;index"`
				TLEID                 string    `gorm:"not null;index"`
				ContextID             string    `gorm:"not null;index"`
				Context               Context   `gorm:"constraint:OnDelete:CASCADE;foreignKey:ContextID;references:ID"`
				Tile                  Tile      `gorm:"constraint:OnDelete:CASCADE;foreignKey:TileID;references:ID"`
				IntersectionLatitude  float64   `gorm:"type:double precision;not null;"`
				IntersectionLongitude float64   `gorm:"type:double precision;not null;"`
				IntersectedAt         time.Time `gorm:"not null"`
			}

			// Define the many-to-many relationship tables
			type ContextSatellite struct {
				ContextID   string    `gorm:"not null;index;uniqueIndex:unique_context_satellite"`
				SatelliteID string    `gorm:"not null;index;uniqueIndex:unique_context_satellite"`
				Context     Context   `gorm:"constraint:OnDelete:CASCADE;foreignKey:ContextID;references:ID"`
				Satellite   Satellite `gorm:"constraint:OnDelete:CASCADE;foreignKey:SatelliteID;references:ID"`
			}

			type ContextTLE struct {
				ContextID string  `gorm:"not null;index;uniqueIndex:unique_context_tle"`
				TLEID     string  `gorm:"not null;index;uniqueIndex:unique_context_tle"`
				Context   Context `gorm:"constraint:OnDelete:CASCADE;foreignKey:ContextID;references:ID"`
				TLE       TLE     `gorm:"constraint:OnDelete:CASCADE;foreignKey:TLEID;references:ID"`
			}

			type ContextTile struct {
				ContextID string  `gorm:"not null;index;uniqueIndex:unique_context_tile"`
				TileID    string  `gorm:"not null;index;uniqueIndex:unique_context_tile"`
				Context   Context `gorm:"constraint:OnDelete:CASCADE;foreignKey:ContextID;references:ID"`
				Tile      Tile    `gorm:"constraint:OnDelete:CASCADE;foreignKey:TileID;references:ID"`
			}

			// AutoMigrate all tables
			if err := db.AutoMigrate(
				&Context{},
				&Satellite{},
				&TLE{},
				&Tile{},
				&TileSatelliteMapping{},
				&ContextSatellite{},
				&ContextTLE{},
				&ContextTile{},
			); err != nil {
				return err
			}

			if err := db.Exec(`
			ALTER TABLE contexts
			ADD CONSTRAINT unique_tenant_context_name
			UNIQUE (tenant_id, name);
		`).Error; err != nil {
				return err
			}

			return nil
		},
		Rollback: func(db *gorm.DB) error {
			// Drop tables in reverse order of dependencies
			if err := db.Exec(`
				ALTER TABLE contexts
				DROP CONSTRAINT IF EXISTS unique_tenant_context_name;
			`).Error; err != nil {
				return err
			}

			return db.Migrator().DropTable(
				"context_satellites",
				"context_tles",
				"context_tiles",
				"tile_satellite_mappings",
				"tiles",
				"tles",
				"satellites",
				"contexts",
			)
		},
	}

	AddMigration(m)
}
