package models

import "time"

type Order struct {
	ID          string
	UserID      string
	ProductsIDs []string
	Status      string
	Total       uint

	CreatedAt time.Time
	UpdatedAt time.Time
}
