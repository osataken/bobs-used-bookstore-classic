package handler

import (
	"net/http"

	"bobs-used-bookstore-api/internal/dto"
	"bobs-used-bookstore-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authMiddleware *middleware.AuthMiddleware
}

func NewAuthHandler(authMiddleware *middleware.AuthMiddleware) *AuthHandler {
	return &AuthHandler{authMiddleware: authMiddleware}
}

func (h *AuthHandler) Login(c *gin.Context) {
	h.authMiddleware.Login(c)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	h.authMiddleware.Logout(c)
}

func (h *AuthHandler) Me(c *gin.Context) {
	sub := middleware.GetUserSub(c)
	if sub == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	username, _ := c.Get(middleware.ContextKeyUsername)
	firstName, _ := c.Get(middleware.ContextKeyFirstName)
	lastName, _ := c.Get(middleware.ContextKeyLastName)
	isAdmin, _ := c.Get(middleware.ContextKeyIsAdmin)

	c.JSON(http.StatusOK, dto.AuthUserResponse{
		Sub:       sub,
		Username:  username.(string),
		FirstName: firstName.(string),
		LastName:  lastName.(string),
		IsAdmin:   isAdmin.(bool),
	})
}
