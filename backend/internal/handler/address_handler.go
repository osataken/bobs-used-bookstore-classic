package handler

import (
	"net/http"
	"strconv"

	"bobs-used-bookstore-api/internal/dto"
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

	var response []dto.AddressResponse
	for _, a := range addresses {
		response = append(response, dto.AddressResponse{
			ID:           a.ID,
			AddressLine1: a.AddressLine1,
			AddressLine2: a.AddressLine2,
			City:         a.City,
			State:        a.State,
			Country:      a.Country,
			ZipCode:      a.ZipCode,
		})
	}

	c.JSON(http.StatusOK, gin.H{"addresses": response})
}

func (h *AddressHandler) Get(c *gin.Context) {
	sub := middleware.GetUserSub(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	address, err := h.addressService.GetAddress(sub, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "address not found"})
		return
	}

	c.JSON(http.StatusOK, dto.AddressResponse{
		ID:           address.ID,
		AddressLine1: address.AddressLine1,
		AddressLine2: address.AddressLine2,
		City:         address.City,
		State:        address.State,
		Country:      address.Country,
		ZipCode:      address.ZipCode,
	})
}

func (h *AddressHandler) Create(c *gin.Context) {
	sub := middleware.GetUserSub(c)

	var req dto.AddressCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.addressService.Create(service.CreateAddressDTO{
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		City:         req.City,
		State:        req.State,
		Country:      req.Country,
		ZipCode:      req.ZipCode,
		CustomerSub:  sub,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create address"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "address created"})
}

func (h *AddressHandler) Update(c *gin.Context) {
	sub := middleware.GetUserSub(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.AddressUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err = h.addressService.Update(service.UpdateAddressDTO{
		AddressID:    id,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		City:         req.City,
		State:        req.State,
		Country:      req.Country,
		ZipCode:      req.ZipCode,
		CustomerSub:  sub,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "address updated"})
}

func (h *AddressHandler) Delete(c *gin.Context) {
	sub := middleware.GetUserSub(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.addressService.Delete(service.DeleteAddressDTO{
		AddressID:   id,
		CustomerSub: sub,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "address deleted"})
}
