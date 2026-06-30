package handler

import (
	"bobs-used-bookstore-api/internal/config"
	"bobs-used-bookstore-api/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	cfg *config.Config
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

func (h *AuthHandler) Login(c *gin.Context) {
	if h.cfg.IsLocalAuth() {
		middleware.LoginLocal(c)
		return
	}
	// Cognito redirect would go here
	c.JSON(http.StatusOK, gin.H{"message": "redirect to cognito"})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	if h.cfg.IsLocalAuth() {
		middleware.LogoutLocal(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *AuthHandler) Me(c *gin.Context) {
	sub, exists := c.Get(middleware.ContextUserSub)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"authenticated": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"authenticated": true,
		"sub":           sub,
		"username":      c.GetString(middleware.ContextUserName),
		"firstName":     c.GetString(middleware.ContextFirstName),
		"lastName":      c.GetString(middleware.ContextLastName),
		"role":          c.GetString(middleware.ContextUserRole),
		"customerID":    c.GetInt(middleware.ContextCustomerID),
	})
}
