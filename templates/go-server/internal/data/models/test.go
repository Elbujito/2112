package models

// Satellite Model
type Test struct {
	ModelBase
	Name string `gorm:"size:255;not null"` // Test name
}

// MapToForm maps the Test model to a TestForm.
func (model *Test) MapToForm() *TestForm {
	return &TestForm{
		FormBase: FormBase{
			ID:        model.ID,
			CreatedAt: model.CreatedAt,
			UpdatedAt: model.UpdatedAt,
		},
		Name: model.Name,
	}
}
