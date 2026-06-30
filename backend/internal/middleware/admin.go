package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminMiddleware struct{}

func NewAdminMiddleware() *AdminMiddleware {
	return &AdminMiddleware{}
}

// RequireAdmin returns middleware that requires administrator role
func (m *AdminMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get(ContextKeyIsAdmin)
		if !exists || !isAdmin.(bool) {
			c.JSON(http.StatusForbidden, gin.H{"error": "administrator access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
