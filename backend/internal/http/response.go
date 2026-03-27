package http

import "pack-calculator/internal/domain"

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status string `json:"status"`
}

// PackSizesResponse represents the response for the pack sizes endpoint.
type PackSizesResponse struct {
	PackSizes []int `json:"pack_sizes"`
}

// CalculateResponse represents the response for the calculate endpoint.
type CalculateResponse struct {
	OrderQuantity int                     `json:"order_quantity"`
	TotalShipped  int                     `json:"total_shipped"`
	TotalPacks    int                     `json:"total_packs"`
	Packs         []domain.PackSelection  `json:"packs"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error string `json:"error"`
}
