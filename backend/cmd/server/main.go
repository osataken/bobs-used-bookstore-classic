package main

import (
	"log/slog"
	"os"

	"bobs-used-bookstore-api/internal/config"
	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/repository"
	"bobs-used-bookstore-api/internal/router"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	db, err := repository.InitDB(cfg)
	if err != nil {
		slog.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}

	repos := repository.NewRepositories(db)
	services := service.NewServices(repos, cfg)
	authMiddleware := middleware.NewAuthMiddleware(cfg)
	adminMiddleware := middleware.NewAdminMiddleware()

	r := gin.Default()
	r.Use(middleware.CORS())

	router.Setup(r, services, authMiddleware, adminMiddleware, cfg)

	slog.Info("starting server", "port", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
