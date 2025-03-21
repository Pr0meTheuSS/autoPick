package models

import "time"

type Product struct {
	ID          string
	Name        string
	Description string
	Price       float32
	Category    string
	Attributes  map[string]any

	CreatedAt time.Time
	UpdatedAt time.Time
}
