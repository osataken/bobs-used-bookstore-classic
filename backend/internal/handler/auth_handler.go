package handler

import (
	"net/http"
	"time"

	"bobs-used-bookstore-api/internal/config"
	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	customerService *service.CustomerService
	cfg             *config.Config
}

func NewAuthHandler(customerService *service.CustomerService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{customerService: customerService, cfg: cfg}
}

func (h *AuthHandler) Login(c *gin.Context) {
	if h.cfg.AuthenticationMode == "aws" {
		// In AWS mode, redirect to Cognito
		c.JSON(http.StatusOK, gin.H{
			"authUrl": h.cfg.CognitoMetadataURL,
			"clientId": h.cfg.CognitoClientID,
		})
		return
	}

	// Local mode: set cookie and create/update customer
	c.SetCookie(
		middleware.LocalAuthCookieName,
		"authenticated",
		int(24*time.Hour/time.Second),
		"/",
		"",
		false,
		false,
	)

	// Create or update local user
	_ = h.customerService.CreateOrUpdate(
		h.cfg.LocalUserSub,
		h.cfg.LocalUsername,
		"Bookstore",
		"User",
	)

	c.JSON(http.StatusOK, gin.H{
		"message":  "logged in",
		"username": h.cfg.LocalUsername,
		"sub":      h.cfg.LocalUserSub,
		"roles":    []string{"Administrators"},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	if h.cfg.AuthenticationMode == "aws" {
		c.JSON(http.StatusOK, gin.H{
			"logoutUrl": h.cfg.CognitoDomain + "/logout",
		})
		return
	}

	c.SetCookie(middleware.LocalAuthCookieName, "", -1, "/", "", false, false)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func (h *AuthHandler) Me(c *gin.Context) {
	sub := middleware.GetUserSub(c)
	if sub == "" {
		c.JSON(http.StatusOK, gin.H{"authenticated": false})
		return
	}

	customer, err := h.customerService.GetBySub(sub)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"authenticated": true,
			"sub":           sub,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"authenticated": true,
		"sub":           customer.Sub,
		"username":      customer.Username,
		"firstName":     customer.FirstName,
		"lastName":      customer.LastName,
		"roles":         middleware.GetUserRoles(c),
	})
}
