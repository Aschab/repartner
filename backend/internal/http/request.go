package http

// CalculateRequest represents the request body for the calculate endpoint.
type CalculateRequest struct {
	OrderQuantity int `json:"order_quantity"`
}
