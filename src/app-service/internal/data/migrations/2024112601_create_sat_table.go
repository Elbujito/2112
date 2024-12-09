package migrations

import (
	"time"

	"github.com/Elbujito/2112/src/app-service/internal/data/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	m := &gormigrate.Migration{
		ID: "2024112604_create_quadkey_visibility_schema",
		Migrate: func(db *gorm.DB) error {
			type Satellite struct {
				models.ModelBase
				Name           string     `gorm:"size:255;not null"`        // Satellite name
				NoradID        string     `gorm:"size:255;unique;not null"` // NORAD ID
				Type           string     `gorm:"size:255"`                 // Satellite type (e.g., telescope, communication)
				LaunchDate     *time.Time `gorm:"type:date"`                // Launch date
				DecayDate      *time.Time `gorm:"type:date"`                // Decay date (optional)
				IntlDesignator string     `gorm:"size:255"`                 // International designator
				Owner          string     `gorm:"size:255"`                 // Ownership information
				ObjectType     string     `gorm:"size:255"`                 // Object type (e.g., "PAYLOAD")
				Period         *float64   `gorm:"type:float"`               // Orbital period in minutes (optional)
				Inclination    *float64   `gorm:"type:float"`               // Orbital inclination in degrees (optional)
				Apogee         *float64   `gorm:"type:float"`               // Apogee altitude in kilometers (optional)
				Perigee        *float64   `gorm:"type:float"`               // Perigee altitude in kilometers (optional)
				RCS            *float64   `gorm:"type:float"`               // Radar cross-section in square meters (optional)
				Altitude       *float64   `gorm:"type:float"`               // Altitude in kilometers (optional)
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
				Quadkey        string  `gorm:"size:256;unique;not null"`                  // Unique identifier for the tile (Quadkey)
				ZoomLevel      int     `gorm:"not null"`                                  // Zoom level for the tile
				CenterLat      float64 `gorm:"not null"`                                  // Center latitude of the tile
				CenterLon      float64 `gorm:"not null"`                                  // Center longitude of the tile
				NbFaces        int     `gorm:"not null"`                                  // Number of faces in the tile's shape
				Radius         float64 `gorm:"not null"`                                  // Radius of the tile in meters
				BoundariesJSON string  `gorm:"type:json"`                                 // Serialized JSON of the boundary vertices of the tile
				SpatialIndex   string  `gorm:"type:geometry(Polygon, 4326);spatialIndex"` // PostGIS geometry type with SRID 4326
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
			if err := db.AutoMigrate(&Satellite{}, &TLE{}, &Tile{}, &TileSatelliteMapping{}); err != nil {
				return err
			}

			// Convert TileSatelliteMapping to a hypertable
			if err := db.Exec(`SELECT create_hypertable('tile_satellite_mappings', 'aos', if_not_exists => TRUE);`).Error; err != nil {
				return err
			}

			// Create a continuous aggregate for 1-hour buckets
			if err := db.Exec(`
				CREATE MATERIALIZED VIEW hourly_visibility
				WITH (timescaledb.continuous) AS
				SELECT
					time_bucket('1 hour', aos) AS bucket,
					norad_id,
					tile_id,
					MAX(max_elevation) AS max_elevation
				FROM
					tile_satellite_mappings
				GROUP BY
					bucket, norad_id, tile_id;
			`).Error; err != nil {
				return err
			}

			// Add retention policy for TileSatelliteMapping
			if err := db.Exec(`
				SELECT add_retention_policy('tile_satellite_mappings', INTERVAL '7 days');
			`).Error; err != nil {
				return err
			}

			// Add retention policy for hourly_visibility (optional)
			if err := db.Exec(`
				SELECT add_retention_policy('hourly_visibility', INTERVAL '7 days');
			`).Error; err != nil {
				return err
			}

			// Ensure existing geometries have SRID 4326
			return db.Exec(`
				UPDATE tiles
				SET spatial_index = ST_SetSRID(spatial_index, 4326)
				WHERE ST_SRID(spatial_index) != 4326;
			`).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Migrator().DropTable("hourly_visibility", "tile_satellite_mappings", "tiles", "tles", "satellites")
		},
	}

	AddMigration(m)
}
