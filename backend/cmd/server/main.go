package main

import (
	"bobs-used-bookstore-api/internal/config"
	"bobs-used-bookstore-api/internal/database"
	"bobs-used-bookstore-api/internal/handler"
	"bobs-used-bookstore-api/internal/repository"
	"bobs-used-bookstore-api/internal/router"
	"bobs-used-bookstore-api/internal/service"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Initialize database
	db, err := database.Initialize(cfg)
	if err != nil {
		slog.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}

	// Repositories
	bookRepo := repository.NewBookRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	addressRepo := repository.NewAddressRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	offerRepo := repository.NewOfferRepository(db)
	cartRepo := repository.NewShoppingCartRepository(db)
	refDataRepo := repository.NewReferenceDataRepository(db)

	// Services
	bookService := service.NewBookService(bookRepo)
	customerService := service.NewCustomerService(customerRepo)
	addressService := service.NewAddressService(addressRepo)
	orderService := service.NewOrderService(orderRepo, bookRepo, cartRepo)
	offerService := service.NewOfferService(offerRepo)
	cartService := service.NewShoppingCartService(cartRepo)
	refDataService := service.NewReferenceDataService(refDataRepo)
	fileService := service.NewFileService(cfg)
	imageService := service.NewImageService(cfg)

	// Handlers
	homeHandler := handler.NewHomeHandler(bookService)
	searchHandler := handler.NewSearchHandler(bookService, cartService)
	cartHandler := handler.NewShoppingCartHandler(cartService)
	wishlistHandler := handler.NewWishlistHandler(cartService)
	checkoutHandler := handler.NewCheckoutHandler(orderService, cartService, addressService)
	orderHandler := handler.NewOrderHandler(orderService)
	addressHandler := handler.NewAddressHandler(addressService)
	resaleHandler := handler.NewResaleHandler(offerService, refDataService)
	authHandler := handler.NewAuthHandler(cfg)
	adminDashboard := handler.NewAdminDashboardHandler(orderService, offerService, bookService)
	adminInventory := handler.NewAdminInventoryHandler(bookService, fileService, imageService)
	adminOrder := handler.NewAdminOrderHandler(orderService)
	adminOffer := handler.NewAdminOfferHandler(offerService)
	adminRefData := handler.NewAdminReferenceDataHandler(refDataService)

	// Setup Gin
	r := gin.Default()

	// CORS for frontend
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	router.Setup(r, cfg, customerService, homeHandler, searchHandler, cartHandler,
		wishlistHandler, checkoutHandler, orderHandler, addressHandler, resaleHandler,
		authHandler, adminDashboard, adminInventory, adminOrder, adminOffer, adminRefData)

	slog.Info("starting server", "port", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
