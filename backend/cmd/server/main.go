package main

import (
	"log/slog"
	"os"

	"bobs-used-bookstore-api/internal/config"
	"bobs-used-bookstore-api/internal/handler"
	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/repository"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	db, err := repository.InitDB(cfg)
	if err != nil {
		slog.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}

	// Repositories
	bookRepo := repository.NewBookRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	cartRepo := repository.NewShoppingCartRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	addressRepo := repository.NewAddressRepository(db)
	offerRepo := repository.NewOfferRepository(db)
	refDataRepo := repository.NewReferenceDataRepository(db)

	// Infrastructure services
	fileService := service.NewFileService(cfg)
	imageResizeService := service.NewImageResizeService()
	imageValidationService := service.NewImageValidationService(cfg)

	// Domain services
	bookService := service.NewBookService(bookRepo, fileService, imageResizeService, imageValidationService)
	orderService := service.NewOrderService(orderRepo, bookRepo, cartRepo, db)
	cartService := service.NewShoppingCartService(cartRepo, bookRepo)
	customerService := service.NewCustomerService(customerRepo)
	addressService := service.NewAddressService(addressRepo)
	offerService := service.NewOfferService(offerRepo)
	refDataService := service.NewReferenceDataService(refDataRepo)

	// Handlers
	homeHandler := handler.NewHomeHandler(bookService)
	searchHandler := handler.NewSearchHandler(bookService, cartService)
	cartHandler := handler.NewShoppingCartHandler(cartService)
	wishlistHandler := handler.NewWishlistHandler(cartService)
	checkoutHandler := handler.NewCheckoutHandler(orderService, cartService, addressService, customerService)
	ordersHandler := handler.NewOrdersHandler(orderService, customerService)
	addressHandler := handler.NewAddressHandler(addressService, customerService)
	resaleHandler := handler.NewResaleHandler(offerService, refDataService, customerService)
	authHandler := handler.NewAuthHandler(cfg, customerService)

	adminDashHandler := handler.NewAdminDashboardHandler(orderService, offerService, bookService)
	adminInventoryHandler := handler.NewAdminInventoryHandler(bookService, refDataService)
	adminOffersHandler := handler.NewAdminOffersHandler(offerService)
	adminOrdersHandler := handler.NewAdminOrdersHandler(orderService)
	adminRefDataHandler := handler.NewAdminReferenceDataHandler(refDataService)

	r := gin.Default()

	// CORS middleware
	r.Use(middleware.CORS())

	// Public routes
	api := r.Group("/api")
	{
		api.GET("/home/featured", homeHandler.GetFeaturedBooks)
		api.GET("/search", searchHandler.Search)
		api.GET("/search/:id", searchHandler.Details)
		api.GET("/cart", cartHandler.GetCart)
		api.POST("/cart/add", cartHandler.AddItem)
		api.POST("/cart/delete", cartHandler.DeleteItem)
		api.GET("/wishlist", wishlistHandler.GetWishlist)
		api.POST("/wishlist/add", wishlistHandler.AddItem)
		api.POST("/wishlist/move", wishlistHandler.MoveToCart)
		api.POST("/wishlist/move-all", wishlistHandler.MoveAllToCart)
		api.POST("/wishlist/delete", wishlistHandler.DeleteItem)
	}

	// Auth routes
	auth := api.Group("/auth")
	{
		auth.GET("/login", authHandler.Login)
		auth.GET("/callback", authHandler.Callback)
		auth.GET("/logout", authHandler.Logout)
		auth.GET("/me", authHandler.Me)
	}

	// Authenticated routes
	authenticated := api.Group("/")
	authenticated.Use(middleware.AuthRequired(cfg, customerService))
	{
		authenticated.GET("/checkout", checkoutHandler.GetCheckout)
		authenticated.POST("/checkout", checkoutHandler.SubmitOrder)
		authenticated.GET("/checkout/finished", checkoutHandler.Finished)
		authenticated.GET("/orders", ordersHandler.ListOrders)
		authenticated.GET("/orders/:id", ordersHandler.GetOrder)
		authenticated.POST("/orders/:id/cancel", ordersHandler.CancelOrder)
		authenticated.GET("/addresses", addressHandler.ListAddresses)
		authenticated.POST("/addresses", addressHandler.CreateAddress)
		authenticated.PUT("/addresses/:id", addressHandler.UpdateAddress)
		authenticated.DELETE("/addresses/:id", addressHandler.DeleteAddress)
		authenticated.GET("/resale", resaleHandler.ListOffers)
		authenticated.POST("/resale", resaleHandler.CreateOffer)
		authenticated.GET("/resale/reference-data", resaleHandler.GetReferenceData)
	}

	// Admin routes
	admin := api.Group("/admin")
	admin.Use(middleware.AuthRequired(cfg, customerService))
	admin.Use(middleware.AdminRequired())
	{
		admin.GET("/dashboard", adminDashHandler.GetDashboard)
		admin.GET("/inventory", adminInventoryHandler.ListBooks)
		admin.GET("/inventory/:id", adminInventoryHandler.GetBook)
		admin.POST("/inventory", adminInventoryHandler.CreateBook)
		admin.PUT("/inventory/:id", adminInventoryHandler.UpdateBook)
		admin.GET("/offers", adminOffersHandler.ListOffers)
		admin.POST("/offers/:id/approve", adminOffersHandler.Approve)
		admin.POST("/offers/:id/reject", adminOffersHandler.Reject)
		admin.POST("/offers/:id/received", adminOffersHandler.Received)
		admin.POST("/offers/:id/paid", adminOffersHandler.Paid)
		admin.GET("/orders", adminOrdersHandler.ListOrders)
		admin.GET("/orders/:id", adminOrdersHandler.GetOrder)
		admin.PUT("/orders/:id/status", adminOrdersHandler.UpdateStatus)
		admin.GET("/reference-data", adminRefDataHandler.List)
		admin.POST("/reference-data", adminRefDataHandler.Create)
		admin.PUT("/reference-data/:id", adminRefDataHandler.Update)
	}

	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	slog.Info("starting server", "port", port)
	if err := r.Run(":" + port); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
