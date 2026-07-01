package router

import (
	"bobs-used-bookstore-api/internal/config"
	"bobs-used-bookstore-api/internal/handler"
	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, services *service.Services, authMW *middleware.AuthMiddleware, adminMW *middleware.AdminMiddleware, cfg *config.Config) {
	// Handlers
	homeHandler := handler.NewHomeHandler(services.Book)
	searchHandler := handler.NewSearchHandler(services.Book)
	cartHandler := handler.NewCartHandler(services.Cart)
	wishlistHandler := handler.NewWishlistHandler(services.Cart)
	checkoutHandler := handler.NewCheckoutHandler(services.Cart, services.Order, services.Address)
	orderHandler := handler.NewOrderHandler(services.Order)
	resaleHandler := handler.NewResaleHandler(services.Offer, services.ReferenceData)
	addressHandler := handler.NewAddressHandler(services.Address)
	authHandler := handler.NewAuthHandler(services.Customer, cfg)
	adminDashboardHandler := handler.NewAdminDashboardHandler(services.Order, services.Offer, services.Book)
	adminInventoryHandler := handler.NewAdminInventoryHandler(services.Book, services.ReferenceData)
	adminOrdersHandler := handler.NewAdminOrdersHandler(services.Order)
	adminOffersHandler := handler.NewAdminOffersHandler(services.Offer)
	adminRefDataHandler := handler.NewAdminReferenceDataHandler(services.ReferenceData)

	api := r.Group("/api")

	// Public routes (with optional auth + cart middleware)
	public := api.Group("")
	public.Use(authMW.OptionalAuth(), middleware.CartMiddleware())
	{
		public.GET("/home", homeHandler.Index)
		public.GET("/search", searchHandler.Index)
		public.GET("/search/:id", searchHandler.Details)
		public.GET("/cart", cartHandler.Index)
		public.POST("/cart/add", cartHandler.AddToCart)
		public.POST("/cart/delete", cartHandler.Delete)
		public.GET("/wishlist", wishlistHandler.Index)
		public.POST("/wishlist/add", wishlistHandler.AddToWishlist)
		public.POST("/wishlist/move", wishlistHandler.MoveToCart)
		public.POST("/wishlist/moveAll", wishlistHandler.MoveAllToCart)
		public.POST("/wishlist/delete", wishlistHandler.Delete)
	}

	// Auth routes
	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/logout", authHandler.Logout)
	api.GET("/auth/me", authMW.OptionalAuth(), authHandler.Me)

	// Authenticated routes
	auth := api.Group("")
	auth.Use(authMW.Authenticate(), middleware.CartMiddleware())
	{
		auth.GET("/checkout", checkoutHandler.Index)
		auth.POST("/checkout", checkoutHandler.CreateOrder)
		auth.GET("/checkout/finished", checkoutHandler.Finished)
		auth.GET("/orders", orderHandler.Index)
		auth.GET("/orders/:id", orderHandler.Details)
		auth.POST("/orders/cancel", orderHandler.Cancel)
		auth.GET("/resale", resaleHandler.Index)
		auth.POST("/resale", resaleHandler.Create)
		auth.GET("/resale/referenceData", resaleHandler.GetReferenceData)
		auth.GET("/addresses", addressHandler.Index)
		auth.GET("/addresses/:id", addressHandler.GetByID)
		auth.POST("/addresses", addressHandler.Create)
		auth.PUT("/addresses", addressHandler.Update)
		auth.POST("/addresses/delete", addressHandler.Delete)
	}

	// Admin routes
	admin := api.Group("/admin")
	admin.Use(authMW.Authenticate(), adminMW.RequireAdmin())
	{
		admin.GET("/dashboard", adminDashboardHandler.Index)
		admin.GET("/inventory", adminInventoryHandler.Index)
		admin.GET("/inventory/:id", adminInventoryHandler.Details)
		admin.POST("/inventory", adminInventoryHandler.Create)
		admin.PUT("/inventory/:id", adminInventoryHandler.Update)
		admin.GET("/orders", adminOrdersHandler.Index)
		admin.GET("/orders/:id", adminOrdersHandler.Details)
		admin.POST("/orders/status", adminOrdersHandler.UpdateStatus)
		admin.GET("/offers", adminOffersHandler.Index)
		admin.POST("/offers/status", adminOffersHandler.UpdateStatus)
		admin.GET("/referenceData", adminRefDataHandler.Index)
		admin.GET("/referenceData/all", adminRefDataHandler.GetAll)
		admin.POST("/referenceData", adminRefDataHandler.Create)
		admin.PUT("/referenceData/:id", adminRefDataHandler.Update)
	}

	// Serve uploaded images
	r.Static("/uploads", "./uploads")
}
