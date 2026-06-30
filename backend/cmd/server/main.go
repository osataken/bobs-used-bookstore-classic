package main

import (
	"log/slog"
	"os"

	"bobs-used-bookstore-api/internal/config"
	"bobs-used-bookstore-api/internal/database"
	"bobs-used-bookstore-api/internal/handler"
	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/repository"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}

	// Initialize repositories
	repos := repository.NewRepositories(db)

	// Initialize services
	services := service.NewServices(repos, cfg)

	// Setup Gin router
	r := gin.Default()

	// Setup middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg, services.Customer)
	adminMiddleware := middleware.NewAdminMiddleware()

	// Setup routes
	handler.SetupRoutes(r, services, authMiddleware, adminMiddleware)

	// Start server
	addr := ":" + cfg.Port
	slog.Info("Starting server", "address", addr)
	if err := r.Run(addr); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
