package models

import (
	"time"
)

type FormBase struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TestForm struct {
	FormBase
	Name string `json:"name" validate:"required,min=2,max=255"` // Satellite name
}

func (form *TestForm) MapToModel() *Test {
	return &Test{
		Name: form.Name,
	}
}
