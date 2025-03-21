package dtos

import "time"

type CreateOrderDto struct {
	ID          string   `json:"id"`
	UserID      string   `json:"user_id"`
	ProductsIDs []string `json:"products_ids"`
	Total       uint     `json:"total"`
}

type OrderDto struct {
	ID          string   `json:"id"`
	UserID      string   `json:"user_id"`
	Status      string   `json:"status"`
	ProductsIDs []string `json:"products_ids"`
	Total       uint     `json:"total"`

	CreatedAt time.Time `json:"created_at"`
}

type UpdateOrderStatusDto struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}
