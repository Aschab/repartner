package app

import (
	"net/http"

	"pack-calculator/internal/config"
	httpHandler "pack-calculator/internal/http"
	"pack-calculator/internal/service"
)

// App represents the application.
type App struct {
	Handler http.Handler
	Config  *config.Config
}

// New creates a new App instance.
func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	calc := service.NewCalculator()
	handler := httpHandler.NewHandler(calc, cfg)

	return &App{
		Handler: handler.SetupRoutes(),
		Config:  cfg,
	}, nil
}
