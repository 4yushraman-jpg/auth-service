package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/4yushraman-jpg/auth-service/internal/config"
	"github.com/4yushraman-jpg/auth-service/internal/database"
	"github.com/4yushraman-jpg/auth-service/internal/handler"
	"github.com/4yushraman-jpg/auth-service/internal/observability"
	"github.com/4yushraman-jpg/auth-service/internal/repository"
	"github.com/4yushraman-jpg/auth-service/internal/routes"
	"github.com/4yushraman-jpg/auth-service/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	logger := observability.NewLogger()

	db, err := database.NewPostgres(cfg)
	if err != nil {
		logger.Error("failed to connect database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	healthHandler := handler.NewHealthHandler()

	userRepo := repository.NewUserRepository(db)

	passwordService := service.NewPasswordService()

	authService := service.NewAuthService(
		userRepo,
		passwordService,
	)

	authHandler := handler.NewAuthHandler(authService)

	router := routes.SetupRouter(
		logger,
		healthHandler,
		authHandler,
	)

	server := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("server started", "port", cfg.HTTPPort)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-stop

	logger.Info("shutting down server")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("shutdown failed", "error", err)
	}

	logger.Info("server exited properly")
}
