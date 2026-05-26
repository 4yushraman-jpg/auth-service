package routes

import (
	"log/slog"

	"github.com/4yushraman-jpg/auth-service/internal/handler"
	"github.com/4yushraman-jpg/auth-service/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func SetupRouter(
	logger *slog.Logger,
	healthHandler *handler.HealthHandler,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logging(logger))

	r.Get("/health", healthHandler.Health)

	return r
}
