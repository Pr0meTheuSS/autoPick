package dtos

import "time"

// DTO для создания продукта
type CreateProductDto struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       float32        `json:"price"`
	Category    string         `json:"category"`
	Attributes  map[string]any `json:"attributes"`
}

// DTO для получения продукта
type ProductDto struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       uint           `json:"price"`
	Category    string         `json:"category"`
	Attributes  map[string]any `json:"attributes"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
