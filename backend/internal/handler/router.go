package handler

import (
	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *gin.Engine, services *service.Services, authMW *middleware.AuthMiddleware, adminMW *middleware.AdminMiddleware) {
	// Initialize handlers
	homeHandler := NewHomeHandler(services.Book)
	searchHandler := NewSearchHandler(services.Book, services.ShoppingCart)
	cartHandler := NewCartHandler(services.ShoppingCart)
	wishlistHandler := NewWishlistHandler(services.ShoppingCart)
	checkoutHandler := NewCheckoutHandler(services.Order, services.ShoppingCart, services.Address)
	orderHandler := NewOrderHandler(services.Order)
	addressHandler := NewAddressHandler(services.Address)
	resaleHandler := NewResaleHandler(services.Offer, services.ReferenceData)
	authHandler := NewAuthHandler(authMW)
	adminDashboard := NewAdminDashboardHandler(services.Order, services.Offer, services.Book)
	adminInventory := NewAdminInventoryHandler(services.Book)
	adminOrders := NewAdminOrderHandler(services.Order)
	adminOffers := NewAdminOfferHandler(services.Offer)
	adminRefData := NewAdminReferenceDataHandler(services.ReferenceData)

	// API routes
	api := r.Group("/api")
	{
		// Public routes (no auth required)
		api.GET("/home", homeHandler.Index)

		// Search (public)
		search := api.Group("/search")
		search.Use(authMW.OptionalAuth())
		{
			search.GET("", searchHandler.Index)
			search.GET("/:id", searchHandler.Details)
			search.POST("/add-to-cart", searchHandler.AddItemToShoppingCart)
			search.POST("/add-to-wishlist", searchHandler.AddItemToWishlist)
		}

		// Shopping Cart (public)
		cart := api.Group("/cart")
		cart.Use(authMW.OptionalAuth())
		{
			cart.GET("", cartHandler.Index)
			cart.POST("/add", cartHandler.AddItem)
			cart.POST("/delete", cartHandler.Delete)
		}

		// Wishlist (public)
		wishlist := api.Group("/wishlist")
		wishlist.Use(authMW.OptionalAuth())
		{
			wishlist.GET("", wishlistHandler.Index)
			wishlist.POST("/move-to-cart", wishlistHandler.MoveToShoppingCart)
			wishlist.POST("/move-all-to-cart", wishlistHandler.MoveAllToShoppingCart)
			wishlist.POST("/delete", wishlistHandler.Delete)
		}

		// Auth routes
		api.GET("/auth/login", authHandler.Login)
		api.GET("/auth/logout", authHandler.Logout)
		api.GET("/auth/me", authMW.RequireAuth(), authHandler.Me)

		// Reference data (public for forms)
		api.GET("/reference-data", adminRefData.GetAll)

		// Protected routes (require auth)
		protected := api.Group("")
		protected.Use(authMW.RequireAuth())
		{
			// Checkout
			protected.GET("/checkout", checkoutHandler.Index)
			protected.POST("/checkout", checkoutHandler.Submit)

			// Orders
			protected.GET("/orders", orderHandler.Index)
			protected.GET("/orders/:id", orderHandler.Details)
			protected.POST("/orders/:id/cancel", orderHandler.Cancel)

			// Addresses
			protected.GET("/addresses", addressHandler.Index)
			protected.GET("/addresses/:id", addressHandler.Get)
			protected.POST("/addresses", addressHandler.Create)
			protected.PUT("/addresses/:id", addressHandler.Update)
			protected.DELETE("/addresses/:id", addressHandler.Delete)

			// Resale offers
			protected.GET("/resale", resaleHandler.Index)
			protected.POST("/resale", resaleHandler.Create)
		}

		// Admin routes (require auth + admin role)
		admin := api.Group("/admin")
		admin.Use(authMW.RequireAuth(), adminMW.RequireAdmin())
		{
			admin.GET("/dashboard", adminDashboard.Index)

			// Inventory
			admin.GET("/inventory", adminInventory.Index)
			admin.GET("/inventory/:id", adminInventory.Details)
			admin.POST("/inventory", adminInventory.Create)
			admin.PUT("/inventory", adminInventory.Update)

			// Orders
			admin.GET("/orders", adminOrders.Index)
			admin.GET("/orders/:id", adminOrders.Details)
			admin.POST("/orders/:id/status", adminOrders.UpdateStatus)

			// Offers
			admin.GET("/offers", adminOffers.Index)
			admin.POST("/offers/:id/approve", adminOffers.Approve)
			admin.POST("/offers/:id/reject", adminOffers.Reject)
			admin.POST("/offers/:id/received", adminOffers.Received)
			admin.POST("/offers/:id/paid", adminOffers.Paid)

			// Reference Data
			admin.GET("/reference-data", adminRefData.Index)
			admin.GET("/reference-data/:id", adminRefData.Get)
			admin.POST("/reference-data", adminRefData.Create)
			admin.PUT("/reference-data/:id", adminRefData.Update)
		}
	}
}
