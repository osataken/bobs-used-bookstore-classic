package middleware

import (
	"net/http"
	"time"

	"bobs-used-bookstore-api/internal/config"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	LocalAuthCookieName = "LocalAuthentication"
	LocalUserSub        = "FB6135C7-1464-4A72-B74E-4B63D343DD09"
	LocalUsername        = "bookstoreuser"

	ContextKeySub       = "user_sub"
	ContextKeyUsername  = "user_username"
	ContextKeyFirstName = "user_first_name"
	ContextKeyLastName  = "user_last_name"
	ContextKeyIsAdmin   = "user_is_admin"
)

type AuthMiddleware struct {
	cfg             *config.Config
	customerService *service.CustomerService
}

func NewAuthMiddleware(cfg *config.Config, customerService *service.CustomerService) *AuthMiddleware {
	return &AuthMiddleware{cfg: cfg, customerService: customerService}
}

// RequireAuth returns middleware that requires authentication
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	if m.cfg.IsAWS("authentication") {
		return m.cognitoAuth()
	}
	return m.localAuth()
}

// OptionalAuth returns middleware that sets user context if authenticated but doesn't require it
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	if m.cfg.IsAWS("authentication") {
		return m.optionalCognitoAuth()
	}
	return m.optionalLocalAuth()
}

func (m *AuthMiddleware) localAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(LocalAuthCookieName)
		if err != nil || cookie == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		m.setLocalUserContext(c)
		m.upsertCustomer(c)
		c.Next()
	}
}

func (m *AuthMiddleware) optionalLocalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(LocalAuthCookieName)
		if err == nil && cookie != "" {
			m.setLocalUserContext(c)
			m.upsertCustomer(c)
		}
		c.Next()
	}
}

func (m *AuthMiddleware) setLocalUserContext(c *gin.Context) {
	c.Set(ContextKeySub, LocalUserSub)
	c.Set(ContextKeyUsername, LocalUsername)
	c.Set(ContextKeyFirstName, "Bookstore")
	c.Set(ContextKeyLastName, "User")
	c.Set(ContextKeyIsAdmin, true)
}

func (m *AuthMiddleware) cognitoAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement full Cognito OIDC token validation
		c.JSON(http.StatusUnauthorized, gin.H{"error": "cognito auth not implemented"})
		c.Abort()
	}
}

func (m *AuthMiddleware) optionalCognitoAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement optional Cognito auth
		c.Next()
	}
}

func (m *AuthMiddleware) upsertCustomer(c *gin.Context) {
	sub, _ := c.Get(ContextKeySub)
	username, _ := c.Get(ContextKeyUsername)
	firstName, _ := c.Get(ContextKeyFirstName)
	lastName, _ := c.Get(ContextKeyLastName)

	_ = m.customerService.CreateOrUpdate(service.CreateOrUpdateCustomerDTO{
		Sub:       sub.(string),
		Username:  username.(string),
		FirstName: firstName.(string),
		LastName:  lastName.(string),
	})
}

// Login handles the local login endpoint
func (m *AuthMiddleware) Login(c *gin.Context) {
	if m.cfg.IsAWS("authentication") {
		// Redirect to Cognito login
		c.JSON(http.StatusOK, gin.H{"redirect": m.cfg.CognitoDomain + "/login"})
		return
	}

	// Set local auth cookie
	c.SetCookie(LocalAuthCookieName, "true", 86400, "/", "", false, true)

	redirectURI := c.Query("redirectUri")
	if redirectURI == "" {
		redirectURI = "/"
	}
	c.JSON(http.StatusOK, gin.H{"redirect": redirectURI})
}

// Logout handles logout
func (m *AuthMiddleware) Logout(c *gin.Context) {
	if m.cfg.IsAWS("authentication") {
		// Clear session and redirect to Cognito logout
		c.SetCookie(LocalAuthCookieName, "", -1, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"redirect": m.cfg.CognitoDomain + "/logout"})
		return
	}

	c.SetCookie(LocalAuthCookieName, "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"redirect": "/"})
}

// GetUserSub extracts the user sub from context
func GetUserSub(c *gin.Context) string {
	sub, exists := c.Get(ContextKeySub)
	if !exists {
		return ""
	}
	return sub.(string)
}

// IsAuthenticated checks if the user is authenticated
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get(ContextKeySub)
	return exists
}

// SetCartCookie sets the shopping cart correlation cookie
func SetCartCookie(c *gin.Context) string {
	sub := GetUserSub(c)
	if sub != "" {
		// Authenticated user: use their sub as correlation ID
		c.SetCookie("ShoppingCartId", sub, 365*24*int(time.Hour/time.Second), "/", "", false, false)
		return sub
	}

	// Anonymous: check for existing cookie or create new
	cartID, err := c.Cookie("ShoppingCartId")
	if err != nil || cartID == "" {
		cartID = generateUUID()
		c.SetCookie("ShoppingCartId", cartID, 365*24*int(time.Hour/time.Second), "/", "", false, false)
	}
	return cartID
}

// GetCartCorrelationID gets the cart correlation ID from cookie
func GetCartCorrelationID(c *gin.Context) string {
	return SetCartCookie(c)
}

func generateUUID() string {
	return uuid.New().String()
}
