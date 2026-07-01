package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminMiddleware struct{}

func NewAdminMiddleware() *AdminMiddleware {
	return &AdminMiddleware{}
}

func (m *AdminMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		roles := GetUserRoles(c)
		for _, role := range roles {
			if role == "Administrators" {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		c.Abort()
	}
}
