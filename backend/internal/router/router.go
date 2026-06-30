package router

import (
	"bobs-used-bookstore-api/internal/config"
	"bobs-used-bookstore-api/internal/handler"
	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

func Setup(
	r *gin.Engine,
	cfg *config.Config,
	customerService *service.CustomerService,
	homeHandler *handler.HomeHandler,
	searchHandler *handler.SearchHandler,
	cartHandler *handler.ShoppingCartHandler,
	wishlistHandler *handler.WishlistHandler,
	checkoutHandler *handler.CheckoutHandler,
	orderHandler *handler.OrderHandler,
	addressHandler *handler.AddressHandler,
	resaleHandler *handler.ResaleHandler,
	authHandler *handler.AuthHandler,
	adminDashboard *handler.AdminDashboardHandler,
	adminInventory *handler.AdminInventoryHandler,
	adminOrder *handler.AdminOrderHandler,
	adminOffer *handler.AdminOfferHandler,
	adminRefData *handler.AdminReferenceDataHandler,
) {
	// Global middleware
	r.Use(middleware.AuthMiddleware(cfg, customerService))
	r.Use(middleware.CartCookieMiddleware())

	// Public API
	api := r.Group("/api")
	{
		api.GET("/home", homeHandler.Index)

		search := api.Group("/search")
		{
			search.GET("", searchHandler.Index)
			search.GET("/:id", searchHandler.Details)
			search.POST("/add-to-cart", searchHandler.AddToCart)
			search.POST("/add-to-wishlist", searchHandler.AddToWishlist)
		}

		cart := api.Group("/cart")
		{
			cart.GET("", cartHandler.Index)
			cart.POST("/delete", cartHandler.Delete)
		}

		wishlist := api.Group("/wishlist")
		{
			wishlist.GET("", wishlistHandler.Index)
			wishlist.POST("/move-to-cart", wishlistHandler.MoveToCart)
			wishlist.POST("/move-all-to-cart", wishlistHandler.MoveAllToCart)
			wishlist.POST("/delete", wishlistHandler.Delete)
		}

		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/me", authHandler.Me)
		}

		api.GET("/reference-data", resaleHandler.GetReferenceData)

		// Authenticated routes
		authenticated := api.Group("")
		authenticated.Use(middleware.RequireAuth())
		{
			checkout := authenticated.Group("/checkout")
			{
				checkout.GET("", checkoutHandler.Index)
				checkout.POST("", checkoutHandler.PlaceOrder)
			}

			orders := authenticated.Group("/orders")
			{
				orders.GET("", orderHandler.Index)
				orders.GET("/:id", orderHandler.Details)
				orders.POST("/cancel", orderHandler.Cancel)
			}

			addresses := authenticated.Group("/addresses")
			{
				addresses.GET("", addressHandler.Index)
				addresses.POST("", addressHandler.Create)
				addresses.PUT("", addressHandler.Update)
				addresses.POST("/delete", addressHandler.Delete)
			}

			resale := authenticated.Group("/resale")
			{
				resale.GET("", resaleHandler.Index)
				resale.POST("", resaleHandler.Create)
			}
		}
	}

	// Admin routes
	admin := r.Group("/api/admin")
	admin.Use(middleware.RequireAuth())
	admin.Use(middleware.AdminMiddleware())
	{
		admin.GET("/dashboard", adminDashboard.Index)

		inventory := admin.Group("/inventory")
		{
			inventory.GET("", adminInventory.Index)
			inventory.GET("/:id", adminInventory.Details)
			inventory.POST("", adminInventory.Create)
			inventory.PUT("", adminInventory.Update)
			inventory.POST("/:id/image", adminInventory.UploadImage)
		}

		adminOrders := admin.Group("/orders")
		{
			adminOrders.GET("", adminOrder.Index)
			adminOrders.GET("/:id", adminOrder.Details)
			adminOrders.POST("/update-status", adminOrder.UpdateStatus)
		}

		offers := admin.Group("/offers")
		{
			offers.GET("", adminOffer.Index)
			offers.POST("/approve", adminOffer.Approve)
			offers.POST("/reject", adminOffer.Reject)
			offers.POST("/received", adminOffer.Received)
			offers.POST("/paid", adminOffer.Paid)
		}

		refData := admin.Group("/reference-data")
		{
			refData.GET("", adminRefData.Index)
			refData.GET("/:id", adminRefData.GetByID)
			refData.POST("", adminRefData.Create)
			refData.PUT("", adminRefData.Update)
		}
	}
}
