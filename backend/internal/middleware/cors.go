package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS middleware
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			origin = "*"
		}
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Cookie")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// ContextKey type for context values
type contextKey string

const (
	UserSubKey   contextKey = "userSub"
	UserNameKey  contextKey = "userName"
	UserRolesKey contextKey = "userRoles"
	CustomerKey  contextKey = "customer"
)

// GetUserSub extracts user sub from gin context
func GetUserSub(c *gin.Context) string {
	if v, exists := c.Get(string(UserSubKey)); exists {
		return v.(string)
	}
	return ""
}

// GetUserRoles extracts user roles from gin context
func GetUserRoles(c *gin.Context) []string {
	if v, exists := c.Get(string(UserRolesKey)); exists {
		return v.([]string)
	}
	return nil
}

// IsAdmin checks if user has Administrators role
func IsAdmin(c *gin.Context) bool {
	roles := GetUserRoles(c)
	for _, role := range roles {
		if strings.EqualFold(role, "Administrators") {
			return true
		}
	}
	return false
}
