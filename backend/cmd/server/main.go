package main

import (
	"bobs-used-bookstore-api/internal/config"
	"bobs-used-bookstore-api/internal/handler"
	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/repository"
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
	db, err := repository.NewDatabase(cfg.DatabaseDSN)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}

	// Repositories
	bookRepo := repository.NewBookRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	addressRepo := repository.NewAddressRepository(db)
	cartRepo := repository.NewShoppingCartRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	offerRepo := repository.NewOfferRepository(db)
	refDataRepo := repository.NewReferenceDataRepository(db)

	// Services
	fileService := service.NewLocalFileService(cfg.UploadDir)
	imageService := service.NewLocalImageValidationService()
	bookService := service.NewBookService(bookRepo, fileService, imageService, cfg.UploadDir)
	customerService := service.NewCustomerService(customerRepo)
	addressService := service.NewAddressService(addressRepo, customerRepo)
	cartService := service.NewShoppingCartService(cartRepo)
	orderService := service.NewOrderService(orderRepo, cartRepo, customerRepo, bookRepo, db)
	offerService := service.NewOfferService(offerRepo, customerRepo)
	refDataService := service.NewReferenceDataService(refDataRepo)

	// Handlers
	bookHandler := handler.NewBookHandler(bookService)
	cartHandler := handler.NewCartHandler(cartService)
	orderHandler := handler.NewOrderHandler(orderService)
	addressHandler := handler.NewAddressHandler(addressService)
	offerHandler := handler.NewOfferHandler(offerService)
	refDataHandler := handler.NewReferenceDataHandler(refDataService)

	// Router setup
	r := gin.Default()

	// Global middleware
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.LocalAuthMiddleware(customerService))
	r.Use(middleware.ShoppingCartMiddleware())

	// Inject cart service into context for book handler
	r.Use(func(c *gin.Context) {
		c.Set("cartService", cartService)
		c.Next()
	})

	// Serve uploaded files
	r.Static("/uploads", cfg.UploadDir)

	// Public API routes
	api := r.Group("/api")
	{
		// Auth
		api.GET("/auth/status", middleware.AuthStatusHandler)
		// Note: /api/auth/login and /api/auth/logout handled by LocalAuthMiddleware

		// Books (public)
		api.GET("/books/search", bookHandler.Search)
		api.GET("/books/best-selling", bookHandler.BestSelling)
		api.GET("/books/:id", bookHandler.GetByID)
		api.POST("/books/:id/add-to-cart", bookHandler.AddToCart)
		api.POST("/books/:id/add-to-wishlist", bookHandler.AddToWishlist)

		// Cart (public)
		api.GET("/cart", cartHandler.GetCart)
		api.DELETE("/cart/:itemId", cartHandler.DeleteItem)

		// Wishlist (public)
		api.GET("/wishlist", cartHandler.GetWishlist)
		api.POST("/wishlist/:itemId/move-to-cart", cartHandler.MoveToCart)
		api.POST("/wishlist/move-all-to-cart", cartHandler.MoveAllToCart)
		api.DELETE("/wishlist/:itemId", cartHandler.DeleteWishlistItem)

		// Reference data (public read)
		api.GET("/reference-data", refDataHandler.GetAll)

		// Authenticated routes
		auth := api.Group("")
		auth.Use(middleware.RequireAuth())
		{
			// Checkout
			auth.POST("/checkout", orderHandler.Checkout)

			// Orders
			auth.GET("/orders", orderHandler.GetMyOrders)
			auth.GET("/orders/:id", orderHandler.GetOrder)
			auth.POST("/orders/:id/cancel", orderHandler.CancelOrder)

			// Addresses
			auth.GET("/addresses", addressHandler.List)
			auth.GET("/addresses/:id", addressHandler.Get)
			auth.POST("/addresses", addressHandler.Create)
			auth.PUT("/addresses/:id", addressHandler.Update)
			auth.DELETE("/addresses/:id", addressHandler.Delete)

			// Offers (customer)
			auth.GET("/offers", offerHandler.GetMyOffers)
			auth.POST("/offers", offerHandler.Create)
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(middleware.RequireAuth())
		admin.Use(middleware.RequireAdmin())
		{
			// Dashboard
			admin.GET("/dashboard", func(c *gin.Context) {
				orderStats, _ := orderService.GetStatistics()
				offerStats, _ := offerService.GetStatistics()
				bookStats, _ := bookService.GetStatistics()
				c.JSON(200, gin.H{
					"orders":    orderStats,
					"offers":    offerStats,
					"inventory": bookStats,
				})
			})

			// Inventory
			admin.GET("/inventory", bookHandler.AdminList)
			admin.GET("/inventory/:id", bookHandler.GetByID)
			admin.POST("/inventory", bookHandler.AdminCreate)
			admin.PUT("/inventory/:id", bookHandler.AdminUpdate)
			admin.GET("/inventory/statistics", bookHandler.AdminStatistics)

			// Orders
			admin.GET("/orders", orderHandler.AdminList)
			admin.GET("/orders/:id", orderHandler.AdminGetOrder)
			admin.PUT("/orders/:id/status", orderHandler.AdminUpdateStatus)
			admin.GET("/orders/statistics", orderHandler.AdminStatistics)

			// Offers
			admin.GET("/offers", offerHandler.AdminList)
			admin.POST("/offers/:id/approve", offerHandler.AdminApprove)
			admin.POST("/offers/:id/reject", offerHandler.AdminReject)
			admin.POST("/offers/:id/received", offerHandler.AdminReceived)
			admin.POST("/offers/:id/paid", offerHandler.AdminPaid)
			admin.GET("/offers/statistics", offerHandler.AdminStatistics)

			// Reference Data
			admin.GET("/reference-data", refDataHandler.AdminList)
			admin.GET("/reference-data/:id", refDataHandler.AdminGet)
			admin.POST("/reference-data", refDataHandler.AdminCreate)
			admin.PUT("/reference-data/:id", refDataHandler.AdminUpdate)
		}
	}

	slog.Info("Starting server", "port", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
