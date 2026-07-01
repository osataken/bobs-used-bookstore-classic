package handler

import (
	"net/http"
	"strconv"

	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type AddressHandler struct {
	addressService *service.AddressService
}

func NewAddressHandler(addressService *service.AddressService) *AddressHandler {
	return &AddressHandler{addressService: addressService}
}

func (h *AddressHandler) Index(c *gin.Context) {
	sub := middleware.GetUserSub(c)

	addresses, err := h.addressService.GetAddresses(sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get addresses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"addresses": addresses})
}

func (h *AddressHandler) GetByID(c *gin.Context) {
	sub := middleware.GetUserSub(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid address id"})
		return
	}

	address, err := h.addressService.GetAddress(sub, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "address not found"})
		return
	}

	c.JSON(http.StatusOK, address)
}

func (h *AddressHandler) Create(c *gin.Context) {
	sub := middleware.GetUserSub(c)

	var req struct {
		AddressLine1 string `json:"addressLine1" binding:"required"`
		AddressLine2 string `json:"addressLine2"`
		City         string `json:"city" binding:"required"`
		State        string `json:"state" binding:"required"`
		Country      string `json:"country" binding:"required"`
		ZipCode      string `json:"zipCode" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.addressService.CreateAddress(sub, service.CreateAddressInput{
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		City:         req.City,
		State:        req.State,
		Country:      req.Country,
		ZipCode:      req.ZipCode,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create address"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "address created"})
}

func (h *AddressHandler) Update(c *gin.Context) {
	sub := middleware.GetUserSub(c)

	var req struct {
		ID           uint   `json:"id" binding:"required"`
		AddressLine1 string `json:"addressLine1" binding:"required"`
		AddressLine2 string `json:"addressLine2"`
		City         string `json:"city" binding:"required"`
		State        string `json:"state" binding:"required"`
		Country      string `json:"country" binding:"required"`
		ZipCode      string `json:"zipCode" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.addressService.UpdateAddress(sub, service.UpdateAddressInput{
		ID:           req.ID,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		City:         req.City,
		State:        req.State,
		Country:      req.Country,
		ZipCode:      req.ZipCode,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "address updated"})
}

func (h *AddressHandler) Delete(c *gin.Context) {
	sub := middleware.GetUserSub(c)

	var req struct {
		ID uint `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	if err := h.addressService.DeleteAddress(sub, req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "address deleted"})
}
