package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const ShoppingCartCookieName = "ShoppingCartId"

func CartCookieMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// If user is authenticated, use their sub as correlation ID
		// Otherwise, use existing cookie or create new UUID
		sub, exists := c.Get(ContextUserSub)
		if exists && sub != "" {
			c.Set("cartCorrelationID", sub.(string))
			c.SetCookie(ShoppingCartCookieName, sub.(string), 86400*30, "/", "", false, false)
		} else {
			cookieVal, err := c.Cookie(ShoppingCartCookieName)
			if err != nil || cookieVal == "" {
				cookieVal = uuid.New().String()
				c.SetCookie(ShoppingCartCookieName, cookieVal, 86400*30, "/", "", false, false)
			}
			c.Set("cartCorrelationID", cookieVal)
		}
		c.Next()
	}
}
