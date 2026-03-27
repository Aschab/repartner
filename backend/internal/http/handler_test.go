package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"pack-calculator/internal/config"
	"pack-calculator/internal/service"
)

func TestHandler_Health(t *testing.T) {
	cfg := &config.Config{PackSizes: []int{250, 500, 1000}}
	calc := service.NewCalculator()
	handler := NewHandler(calc, cfg)
	router := handler.SetupRoutes()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp HealthResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Status != "ok" {
		t.Errorf("expected status 'ok', got '%s'", resp.Status)
	}
}

func TestHandler_GetPacks(t *testing.T) {
	cfg := &config.Config{PackSizes: []int{250, 500, 1000, 2000, 5000}}
	calc := service.NewCalculator()
	handler := NewHandler(calc, cfg)
	router := handler.SetupRoutes()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/packs", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp PackSizesResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(resp.PackSizes) != 5 {
		t.Errorf("expected 5 pack sizes, got %d", len(resp.PackSizes))
	}
}

func TestHandler_Calculate_ValidRequest(t *testing.T) {
	cfg := &config.Config{PackSizes: []int{250, 500, 1000, 2000, 5000}}
	calc := service.NewCalculator()
	handler := NewHandler(calc, cfg)
	router := handler.SetupRoutes()

	body := CalculateRequest{OrderQuantity: 501}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp CalculateResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.OrderQuantity != 501 {
		t.Errorf("expected order_quantity 501, got %d", resp.OrderQuantity)
	}

	if resp.TotalShipped != 750 {
		t.Errorf("expected total_shipped 750, got %d", resp.TotalShipped)
	}

	if resp.TotalPacks != 2 {
		t.Errorf("expected total_packs 2, got %d", resp.TotalPacks)
	}
}

func TestHandler_Calculate_InvalidRequest_ZeroQuantity(t *testing.T) {
	cfg := &config.Config{PackSizes: []int{250, 500, 1000}}
	calc := service.NewCalculator()
	handler := NewHandler(calc, cfg)
	router := handler.SetupRoutes()

	body := CalculateRequest{OrderQuantity: 0}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var resp ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Error == "" {
		t.Error("expected error message, got empty string")
	}
}

func TestHandler_Calculate_InvalidRequest_NegativeQuantity(t *testing.T) {
	cfg := &config.Config{PackSizes: []int{250, 500, 1000}}
	calc := service.NewCalculator()
	handler := NewHandler(calc, cfg)
	router := handler.SetupRoutes()

	body := CalculateRequest{OrderQuantity: -1}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandler_Calculate_InvalidJSON(t *testing.T) {
	cfg := &config.Config{PackSizes: []int{250, 500, 1000}}
	calc := service.NewCalculator()
	handler := NewHandler(calc, cfg)
	router := handler.SetupRoutes()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandler_CORS(t *testing.T) {
	cfg := &config.Config{PackSizes: []int{250, 500, 1000}}
	calc := service.NewCalculator()
	handler := NewHandler(calc, cfg)
	router := handler.SetupRoutes()

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/calculate", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("expected CORS header to be set")
	}
}
