package middleware

import (
	"net/http"
	"strings"
	"time"

	"bobs-used-bookstore-api/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	LocalAuthCookieName   = "LocalAuthentication"
	ShoppingCartCookieName = "ShoppingCartId"
	ContextKeySub         = "userSub"
	ContextKeyUsername    = "username"
	ContextKeyRoles      = "userRoles"
	ContextKeyCartID     = "shoppingCartId"
)

type AuthMiddleware struct {
	cfg *config.Config
}

func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{cfg: cfg}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	if m.cfg.AuthenticationMode == "aws" {
		return m.cognitoAuth()
	}
	return m.localAuth()
}

func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	if m.cfg.AuthenticationMode == "aws" {
		return m.optionalCognitoAuth()
	}
	return m.optionalLocalAuth()
}

func (m *AuthMiddleware) localAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(LocalAuthCookieName)
		if err != nil || cookie == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Set(ContextKeySub, m.cfg.LocalUserSub)
		c.Set(ContextKeyUsername, m.cfg.LocalUsername)
		c.Set(ContextKeyRoles, []string{"Administrators"})
		c.Next()
	}
}

func (m *AuthMiddleware) optionalLocalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(LocalAuthCookieName)
		if err == nil && cookie != "" {
			c.Set(ContextKeySub, m.cfg.LocalUserSub)
			c.Set(ContextKeyUsername, m.cfg.LocalUsername)
			c.Set(ContextKeyRoles, []string{"Administrators"})
		}
		c.Next()
	}
}

func (m *AuthMiddleware) cognitoAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			c.Abort()
			return
		}

		// In AWS mode, validate JWT token with Cognito
		// For now, extract sub from token (placeholder)
		c.Set(ContextKeySub, token)
		c.Next()
	}
}

func (m *AuthMiddleware) optionalCognitoAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token != authHeader {
				c.Set(ContextKeySub, token)
			}
		}
		c.Next()
	}
}

func CartMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sub, exists := c.Get(ContextKeySub)
		if exists && sub != "" {
			c.Set(ContextKeyCartID, sub.(string))
		} else {
			cartID, err := c.Cookie(ShoppingCartCookieName)
			if err != nil || cartID == "" {
				cartID = uuid.New().String()
			}
			c.Set(ContextKeyCartID, cartID)
			c.SetCookie(ShoppingCartCookieName, cartID, int(24*time.Hour/time.Second), "/", "", false, false)
		}
		c.Next()
	}
}

func GetUserSub(c *gin.Context) string {
	sub, _ := c.Get(ContextKeySub)
	if sub == nil {
		return ""
	}
	return sub.(string)
}

func GetCartID(c *gin.Context) string {
	cartID, _ := c.Get(ContextKeyCartID)
	if cartID == nil {
		return ""
	}
	return cartID.(string)
}

func GetUserRoles(c *gin.Context) []string {
	roles, _ := c.Get(ContextKeyRoles)
	if roles == nil {
		return nil
	}
	return roles.([]string)
}
