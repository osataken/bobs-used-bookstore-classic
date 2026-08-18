package handler

import (
	"net/http"
	"strconv"

	"bobs-used-bookstore-api/internal/middleware"
	"bobs-used-bookstore-api/internal/model"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type AddressHandler struct {
	addressService  *service.AddressService
	customerService *service.CustomerService
}

func NewAddressHandler(addressService *service.AddressService, customerService *service.CustomerService) *AddressHandler {
	return &AddressHandler{addressService: addressService, customerService: customerService}
}

func (h *AddressHandler) ListAddresses(c *gin.Context) {
	customer := c.MustGet(string(middleware.CustomerKey)).(*model.Customer)

	addresses, err := h.addressService.ListByCustomer(customer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch addresses"})
		return
	}
	c.JSON(http.StatusOK, addresses)
}

func (h *AddressHandler) CreateAddress(c *gin.Context) {
	var req struct {
		AddressLine1 string `json:"addressLine1" binding:"required"`
		AddressLine2 string `json:"addressLine2"`
		City         string `json:"city" binding:"required"`
		State        string `json:"state" binding:"required"`
		Country      string `json:"country" binding:"required"`
		ZipCode      string `json:"zipCode" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	customer := c.MustGet(string(middleware.CustomerKey)).(*model.Customer)

	address := &model.Address{
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		City:         req.City,
		State:        req.State,
		Country:      req.Country,
		ZipCode:      req.ZipCode,
		CustomerID:   customer.ID,
		IsActive:     true,
	}

	if err := h.addressService.CreateAddress(address); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create address"})
		return
	}
	c.JSON(http.StatusCreated, address)
}

func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid address id"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	address, err := h.addressService.GetAddress(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "address not found"})
		return
	}

	// Verify ownership
	customer := c.MustGet(string(middleware.CustomerKey)).(*model.Customer)
	if address.CustomerID != customer.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	address.AddressLine1 = req.AddressLine1
	address.AddressLine2 = req.AddressLine2
	address.City = req.City
	address.State = req.State
	address.Country = req.Country
	address.ZipCode = req.ZipCode

	if err := h.addressService.UpdateAddress(address); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update address"})
		return
	}
	c.JSON(http.StatusOK, address)
}

func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid address id"})
		return
	}

	// Verify ownership
	address, err := h.addressService.GetAddress(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "address not found"})
		return
	}
	customer := c.MustGet(string(middleware.CustomerKey)).(*model.Customer)
	if address.CustomerID != customer.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.addressService.DeleteAddress(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete address"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "address deleted"})
}
