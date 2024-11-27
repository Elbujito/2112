package models

// TileService defines the operations required for Tile management
type TileService interface {
	FindByQuadkey(quadkey string) (*Tile, error)
	FindAll() ([]*Tile, error)
	Save(tile *Tile) error
	Create() error
	Update(tile *Tile) error
	Delete(id string) error
}

var tile *Tile = &Tile{}

// Tile Model
type Tile struct {
	ModelBase
	Quadkey   string  `gorm:"size:256;unique;not null"` // Unique identifier for the tile (Quadkey)
	ZoomLevel int     `gorm:"not null"`                 // Zoom level for the tile
	CenterLat float64 `gorm:"not null"`                 // Center latitude of the tile
	CenterLon float64 `gorm:"not null"`                 // Center longitude of the tile
}

func TileModel() *Tile {
	return tile
}

func (model *Tile) FindAll() (models []*Tile, err error) {
	result := db.Model(model).Find(&models)
	return models, result.Error
}

func (model *Tile) FindByQuadkey(quadkey string) (m *Tile, err error) {
	result := db.Model(model).Where("quadkey=?", quadkey).First(&m)
	return m, result.Error
}

func (model *Tile) Save(m *Tile) error {
	return db.Model(model).Create(&m).Error
}

func (model *Tile) Create() error {
	return db.Create(model).Error
}

func (model *Tile) Update(m *Tile) error {
	return db.Model(model).Save(m).Error
}

func (model *Tile) Delete(id string) error {
	return db.Model(model).Where("ID=?", id).Delete(&model).Error
}

func (model *Tile) MapToForm() *TileForm {
	return &TileForm{
		FormBase: FormBase{
			ID:        model.ID,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		},
		Quadkey:   model.Quadkey,
		ZoomLevel: model.ZoomLevel,
		CenterLat: model.CenterLat,
		CenterLon: model.CenterLon,
	}
}
