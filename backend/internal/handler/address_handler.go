package handler

import (
	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AddressHandler struct {
	addressService *service.AddressService
}

func NewAddressHandler(addressService *service.AddressService) *AddressHandler {
	return &AddressHandler{addressService: addressService}
}

// GET /api/addresses
func (h *AddressHandler) List(c *gin.Context) {
	sub := middleware.GetUserSub(c)
	addresses, err := h.addressService.GetAddresses(sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, addresses)
}

// GET /api/addresses/:id
func (h *AddressHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	sub := middleware.GetUserSub(c)
	address, err := h.addressService.GetAddress(sub, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if address == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "address not found"})
		return
	}
	c.JSON(http.StatusOK, address)
}

// POST /api/addresses
func (h *AddressHandler) Create(c *gin.Context) {
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

	sub := middleware.GetUserSub(c)
	dto := service.CreateAddressDTO{
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		City:         req.City,
		State:        req.State,
		Country:      req.Country,
		ZipCode:      req.ZipCode,
	}

	err := h.addressService.Create(sub, dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Address created"})
}

// PUT /api/addresses/:id
func (h *AddressHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

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

	sub := middleware.GetUserSub(c)
	dto := service.UpdateAddressDTO{
		ID:           id,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		City:         req.City,
		State:        req.State,
		Country:      req.Country,
		ZipCode:      req.ZipCode,
	}

	err = h.addressService.Update(sub, dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Address updated"})
}

// DELETE /api/addresses/:id
func (h *AddressHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	sub := middleware.GetUserSub(c)
	err = h.addressService.Delete(sub, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Address deleted"})
}
