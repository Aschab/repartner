package domain

// PackSelection represents a pack size and the number of packs of that size.
type PackSelection struct {
	PackSize int `json:"pack_size"`
	Count    int `json:"count"`
}

// CalculationResult represents the result of a pack calculation.
type CalculationResult struct {
	RequestedQty int             `json:"order_quantity"`
	TotalShipped int             `json:"total_shipped"`
	TotalPacks   int             `json:"total_packs"`
	Packs        []PackSelection `json:"packs"`
}
