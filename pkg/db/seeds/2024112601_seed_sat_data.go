package seeds

import (
	"time"

	"github.com/Elbujito/2112/pkg/db/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	var s = &gormigrate.Migration{}
	s.ID = "2024112602_seed_satellite_and_tle_data"

	s.Migrate = func(db *gorm.DB) error {
		// Seed Satellite Data
		satellites := []*models.Satellite{
			{
				Name:    "Hubble Space Telescope",
				NoradID: "20580",
				Type:    "Telescope",
			},
			{
				Name:    "International Space Station",
				NoradID: "25544",
				Type:    "Space Station",
			},
			{
				Name:    "GPS IIR-10",
				NoradID: "26605",
				Type:    "Navigation",
			},
			{
				Name:    "Landsat 8",
				NoradID: "39084",
				Type:    "Earth Observation",
			},
			{
				Name:    "Starlink-1500",
				NoradID: "45657",
				Type:    "Communication",
			},
		}

		for _, satellite := range satellites {
			if err := satellite.Create(); err != nil {
				logFail(s.ID, err)
				return err
			}
		}

		// Seed TLE Data
		tles := []*models.TLE{
			{
				NoradID: "20580",
				Line1:   "1 20580U 90037B   20245.18473241  .00000238  00000-0  13644-4 0  9996",
				Line2:   "2 20580  28.4708 354.9314 0002627  97.8557 262.2526 14.76817360  4454",
				Epoch:   time.Now().Add(-24 * time.Hour),
			},
			{
				NoradID: "25544",
				Line1:   "1 25544U 98067A   21273.75450833  .00001264  00000-0  29647-4 0  9998",
				Line2:   "2 25544  51.6441 245.2066 0003157  97.5202 262.6127 15.48907224281129",
				Epoch:   time.Now(),
			},
		}

		for _, tle := range tles {
			if err := tle.Create(); err != nil {
				logFail(s.ID, err)
				return err
			}
		}

		// Seed Tile Data
		tiles := []*models.Tile{
			{
				Quadkey:   "213",
				ZoomLevel: 3,
				CenterLat: 52.52,
				CenterLon: 13.405,
			},
			{
				Quadkey:   "123",
				ZoomLevel: 3,
				CenterLat: 40.7128,
				CenterLon: -74.006,
			},
		}

		for _, tile := range tiles {
			if err := tile.Create(); err != nil {
				logFail(s.ID, err)
				return err
			}
		}

		// Seed Visibility Data
		visibilities := []*models.Visibility{
			{
				NoradID:      "20580",
				TileID:       "1",
				StartTime:    time.Now(),
				EndTime:      time.Now().Add(15 * time.Minute),
				MaxElevation: 45.0,
			},
			{
				NoradID:      "25544",
				TileID:       "2",
				StartTime:    time.Now().Add(30 * time.Minute),
				EndTime:      time.Now().Add(45 * time.Minute),
				MaxElevation: 60.0,
			},
		}

		for _, visibility := range visibilities {
			if err := visibility.Save(); err != nil {
				logFail(s.ID, err)
				return err
			}
		}

		logSuccess(s.ID)
		return nil
	}

	AddSeed(s)
}
