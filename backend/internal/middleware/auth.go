package middleware

import (
	"bobs-used-bookstore-api/internal/config"
	"bobs-used-bookstore-api/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	LocalAuthCookie = "LocalAuthentication"
	ContextUserSub  = "userSub"
	ContextUserName = "userName"
	ContextFirstName = "firstName"
	ContextLastName  = "lastName"
	ContextUserRole  = "userRole"
	ContextCustomerID = "customerID"
)

func AuthMiddleware(cfg *config.Config, customerService *service.CustomerService) gin.HandlerFunc {
	if cfg.IsLocalAuth() {
		return localAuthMiddleware(customerService)
	}
	return cognitoAuthMiddleware(cfg, customerService)
}

func localAuthMiddleware(customerService *service.CustomerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(LocalAuthCookie)
		if err != nil || cookie == "" {
			c.Next()
			return
		}

		// Parse the local auth cookie (it's a simple JWT-like token with claims)
		token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
			return []byte("local-dev-secret-key"), nil
		})
		if err != nil || !token.Valid {
			c.Next()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Next()
			return
		}

		sub, _ := claims["sub"].(string)
		username, _ := claims["username"].(string)
		firstName, _ := claims["firstName"].(string)
		lastName, _ := claims["lastName"].(string)
		role, _ := claims["role"].(string)

		// CreateOrUpdateCustomer on every authenticated request
		customer, err := customerService.CreateOrUpdate(sub, username, firstName, lastName)
		if err != nil {
			c.Next()
			return
		}

		c.Set(ContextUserSub, sub)
		c.Set(ContextUserName, username)
		c.Set(ContextFirstName, firstName)
		c.Set(ContextLastName, lastName)
		c.Set(ContextUserRole, role)
		c.Set(ContextCustomerID, customer.ID)
		c.Next()
	}
}

func cognitoAuthMiddleware(cfg *config.Config, customerService *service.CustomerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// In AWS mode, validate Bearer token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenStr := authHeader[7:]
			// In production, validate against Cognito JWKS
			// For now, parse without validation for local testing
			token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
			if err != nil {
				c.Next()
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				c.Next()
				return
			}

			sub, _ := claims["sub"].(string)
			username, _ := claims["cognito:username"].(string)
			firstName, _ := claims["given_name"].(string)
			lastName, _ := claims["family_name"].(string)

			groups, _ := claims["cognito:groups"].([]interface{})
			role := ""
			for _, g := range groups {
				if gs, ok := g.(string); ok && gs == "Administrators" {
					role = "Administrators"
					break
				}
			}

			customer, err := customerService.CreateOrUpdate(sub, username, firstName, lastName)
			if err != nil {
				c.Next()
				return
			}

			c.Set(ContextUserSub, sub)
			c.Set(ContextUserName, username)
			c.Set(ContextFirstName, firstName)
			c.Set(ContextLastName, lastName)
			c.Set(ContextUserRole, role)
			c.Set(ContextCustomerID, customer.ID)
		}
		c.Next()
	}
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get(ContextUserSub)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func LoginLocal(c *gin.Context) {
	// Create a local auth token for development
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":       "local-user-sub-001",
		"username":  "localuser",
		"firstName": "Local",
		"lastName":  "User",
		"role":      "Administrators",
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte("local-dev-secret-key"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	c.SetCookie(LocalAuthCookie, tokenString, 86400, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message":   "logged in",
		"sub":       "local-user-sub-001",
		"username":  "localuser",
		"firstName": "Local",
		"lastName":  "User",
		"role":      "Administrators",
	})
}

func LogoutLocal(c *gin.Context) {
	c.SetCookie(LocalAuthCookie, "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
