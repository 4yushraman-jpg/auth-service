package routes

import (
	"log/slog"

	"github.com/go-chi/chi/v5"

	"github.com/4yushraman-jpg/auth-service/internal/handler"
	"github.com/4yushraman-jpg/auth-service/internal/middleware"
)

func SetupRouter(
	logger *slog.Logger,
	healthHandler *handler.HealthHandler,
	authHandler *handler.AuthHandler,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logging(logger))

	r.Get("/health", healthHandler.Health)

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
	})

	return r
}
