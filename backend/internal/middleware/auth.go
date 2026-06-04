package middleware

import (
	"bobs-used-bookstore-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	LocalUserSub       = "FB6135C7-1464-4A72-B74E-4B63D343DD09"
	LocalUsername       = "bookstoreuser"
	LocalFirstName     = "Bookstore"
	LocalLastName      = "User"
	LocalRole          = "Administrators"
	LocalAuthCookie    = "LocalAuthentication"
	ShoppingCartCookie = "ShoppingCartId"
)

// ContextKey constants
const (
	ContextKeySub  = "userSub"
	ContextKeyRole = "userRole"
	ContextKeyName = "userName"
)

// LocalAuthMiddleware handles local authentication mode
func LocalAuthMiddleware(customerService *service.CustomerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for login endpoint
		if c.Request.URL.Path == "/api/auth/login" {
			// Set auth cookie and create customer
			c.SetCookie(LocalAuthCookie, "true", 86400, "/", "", false, true)
			_ = customerService.CreateOrUpdate(LocalUserSub, LocalUsername, LocalFirstName, LocalLastName)
			c.JSON(http.StatusOK, gin.H{
				"sub":      LocalUserSub,
				"username": LocalUsername,
				"role":     LocalRole,
			})
			c.Abort()
			return
		}

		// Check for logout endpoint
		if c.Request.URL.Path == "/api/auth/logout" {
			c.SetCookie(LocalAuthCookie, "", -1, "/", "", false, true)
			c.JSON(http.StatusOK, gin.H{"message": "logged out"})
			c.Abort()
			return
		}

		// Check if user has local auth cookie
		cookie, err := c.Cookie(LocalAuthCookie)
		if err == nil && cookie != "" {
			_ = customerService.CreateOrUpdate(LocalUserSub, LocalUsername, LocalFirstName, LocalLastName)
			c.Set(ContextKeySub, LocalUserSub)
			c.Set(ContextKeyRole, LocalRole)
			c.Set(ContextKeyName, LocalUsername)
		}

		c.Next()
	}
}

// RequireAuth middleware ensures the user is authenticated
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sub, exists := c.Get(ContextKeySub)
		if !exists || sub == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireAdmin middleware ensures the user has Administrators role
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(ContextKeyRole)
		if !exists || role != LocalRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// ShoppingCartMiddleware manages the shopping cart cookie
func ShoppingCartMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cartID, err := c.Cookie(ShoppingCartCookie)
		if err != nil || cartID == "" {
			sub, exists := c.Get(ContextKeySub)
			if exists && sub != "" {
				cartID = sub.(string)
			} else {
				cartID = uuid.New().String()
			}
		}
		c.SetCookie(ShoppingCartCookie, cartID, 365*24*3600, "/", "", false, false)
		c.Set("shoppingCartId", cartID)
		c.Next()
	}
}

// GetShoppingCartID extracts the shopping cart correlation ID from context
func GetShoppingCartID(c *gin.Context) string {
	id, _ := c.Get("shoppingCartId")
	if id == nil {
		return ""
	}
	return id.(string)
}

// GetUserSub extracts the user sub from context
func GetUserSub(c *gin.Context) string {
	sub, _ := c.Get(ContextKeySub)
	if sub == nil {
		return ""
	}
	return sub.(string)
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// AuthStatusMiddleware returns current auth status
func AuthStatusHandler(c *gin.Context) {
	sub, exists := c.Get(ContextKeySub)
	if !exists || sub == "" {
		c.JSON(http.StatusOK, gin.H{"authenticated": false})
		return
	}

	role, _ := c.Get(ContextKeyRole)
	name, _ := c.Get(ContextKeyName)

	c.JSON(http.StatusOK, gin.H{
		"authenticated": true,
		"sub":           sub,
		"role":          role,
		"username":      name,
	})
}

// Placeholder removed - no longer needed
