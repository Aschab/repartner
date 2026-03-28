package http

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"pack-calculator/internal/config"
	"pack-calculator/internal/service"
)

// Handler handles HTTP requests.
type Handler struct {
	calculator service.Calculator
	config     *config.Config
}

// NewHandler creates a new Handler instance.
func NewHandler(calc service.Calculator, cfg *config.Config) *Handler {
	return &Handler{
		calculator: calc,
		config:     cfg,
	}
}

// SetupRoutes configures the HTTP routes.
func (h *Handler) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("GET /health", h.handleHealth)
	mux.HandleFunc("GET /api/v1/packs", h.handleGetPacks)
	mux.HandleFunc("POST /api/v1/calculate", h.handleCalculate)

	// Serve static files (frontend) - check if static dir exists
	staticDir := "static"
	if _, err := os.Stat(staticDir); err == nil {
		fs := http.FileServer(http.Dir(staticDir))
		mux.Handle("/", h.spaHandler(fs, staticDir))
	}

	return h.corsMiddleware(h.loggingMiddleware(mux))
}

// spaHandler wraps file server to handle SPA routing (serve index.html for unknown routes)
func (h *Handler) spaHandler(fs http.Handler, staticDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := staticDir + r.URL.Path
		// Check if file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// Serve index.html for SPA routing
			http.ServeFile(w, r, staticDir+"/index.html")
			return
		}
		fs.ServeHTTP(w, r)
	})
}

func (h *Handler) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusOK, HealthResponse{Status: "ok"})
}

func (h *Handler) handleGetPacks(w http.ResponseWriter, r *http.Request) {
	h.writeJSON(w, http.StatusOK, PackSizesResponse{PackSizes: h.config.PackSizes})
}

func (h *Handler) handleCalculate(w http.ResponseWriter, r *http.Request) {
	var req CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON body"})
		return
	}

	// Use pack sizes from request, fall back to config if not provided
	packSizes := req.PackSizes
	if len(packSizes) == 0 {
		packSizes = h.config.PackSizes
	}

	result, err := h.calculator.Calculate(req.OrderQuantity, packSizes)
	if err != nil {
		h.writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	resp := CalculateResponse{
		OrderQuantity: result.RequestedQty,
		TotalShipped:  result.TotalShipped,
		TotalPacks:    result.TotalPacks,
		Packs:         result.Packs,
	}

	h.writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}
