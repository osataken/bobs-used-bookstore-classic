package handler

import (
	"net/http"
	"time"

	"bobs-used-bookstore-api/internal/config"
	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	cfg             *config.Config
	customerService *service.CustomerService
}

func NewAuthHandler(cfg *config.Config, customerService *service.CustomerService) *AuthHandler {
	return &AuthHandler{cfg: cfg, customerService: customerService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	if h.cfg.AuthenticationMode == "local" {
		// Generate local auth token
		token, expiry, err := middleware.GenerateLocalToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}

		// Set cookie (1-day expiry, matching source)
		c.SetCookie(middleware.LocalAuthCookieName, token, int(time.Until(expiry).Seconds()), "/", "", false, true)

		// Ensure customer exists
		_, err = h.customerService.CreateOrUpdate(
			middleware.LocalUserSub,
			middleware.LocalUsername,
			middleware.LocalUsername,
			"",
			"",
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "logged in",
			"token":   token,
		})
	} else {
		// Redirect to Cognito
		authURL := h.cfg.CognitoDomain + "/login?" +
			"client_id=" + h.cfg.CognitoClientID +
			"&response_type=code" +
			"&scope=openid+profile+email" +
			"&redirect_uri=" + h.cfg.CognitoRedirectURL
		c.JSON(http.StatusOK, gin.H{"redirectUrl": authURL})
	}
}

func (h *AuthHandler) Callback(c *gin.Context) {
	// In Cognito mode, exchange code for tokens
	// This is a stub - full implementation would use OIDC code exchange
	c.JSON(http.StatusOK, gin.H{"message": "callback received"})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	if h.cfg.AuthenticationMode == "local" {
		c.SetCookie(middleware.LocalAuthCookieName, "", -1, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "logged out"})
	} else {
		logoutURL := h.cfg.CognitoDomain + "/logout?" +
			"client_id=" + h.cfg.CognitoClientID +
			"&logout_uri=" + h.cfg.FrontendURL
		c.JSON(http.StatusOK, gin.H{"redirectUrl": logoutURL})
	}
}

func (h *AuthHandler) Me(c *gin.Context) {
	// Try to get user info from auth
	sub := middleware.GetUserSub(c)
	if sub == "" {
		// Try local auth
		cookie, err := c.Cookie(middleware.LocalAuthCookieName)
		if err != nil || cookie == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}
	}

	customer, exists := c.Get(string(middleware.CustomerKey))
	if !exists {
		// Try to fetch by sub
		sub = middleware.GetUserSub(c)
		cust, err := h.customerService.GetBySub(sub)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
			return
		}
		customer = cust
	}

	cust := customer.(*model.Customer)
	roles := middleware.GetUserRoles(c)

	c.JSON(http.StatusOK, gin.H{
		"id":        cust.ID,
		"sub":       cust.Sub,
		"username":  cust.Username,
		"firstName": cust.FirstName,
		"lastName":  cust.LastName,
		"email":     cust.Email,
		"isAdmin":   containsRole(roles, "Administrators"),
	})
}

func containsRole(roles []string, target string) bool {
	for _, r := range roles {
		if r == target {
			return true
		}
	}
	return false
}
