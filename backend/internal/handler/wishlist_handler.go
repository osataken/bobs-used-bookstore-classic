package handler

import (
	"net/http"

	"bobs-used-bookstore-api/internal/dto"
	"bobs-used-bookstore-api/internal/service"

	"github.com/gin-gonic/gin"
)

type WishlistHandler struct {
	cartService *service.ShoppingCartService
}

func NewWishlistHandler(cartService *service.ShoppingCartService) *WishlistHandler {
	return &WishlistHandler{cartService: cartService}
}

func (h *WishlistHandler) Index(c *gin.Context) {
	correlationID := getCartCorrelationID(c)

	cart, err := h.cartService.GetShoppingCart(correlationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get wishlist"})
		return
	}

	wishlistItems := cart.GetWishlistItems()
	var items []dto.CartItemResponse
	for _, item := range wishlistItems {
		items = append(items, dto.CartItemResponse{
			ID:             item.ID,
			ShoppingCartID: item.ShoppingCartID,
			BookID:         item.BookID,
			Quantity:       item.Quantity,
			WantToBuy:      item.WantToBuy,
			Book:           mapBookToResponse(&item.Book),
		})
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *WishlistHandler) MoveToShoppingCart(c *gin.Context) {
	var req dto.MoveToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	correlationID := getCartCorrelationID(c)

	err := h.cartService.MoveWishlistItemToShoppingCart(service.MoveWishlistItemDTO{
		CorrelationID:      correlationID,
		ShoppingCartItemID: req.ShoppingCartItemID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to move item to cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item moved to cart"})
}

func (h *WishlistHandler) MoveAllToShoppingCart(c *gin.Context) {
	correlationID := getCartCorrelationID(c)

	err := h.cartService.MoveAllWishlistItemsToShoppingCart(service.MoveAllWishlistItemsDTO{
		CorrelationID: correlationID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to move items to cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "all items moved to cart"})
}

func (h *WishlistHandler) Delete(c *gin.Context) {
	var req dto.DeleteCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	correlationID := getCartCorrelationID(c)

	err := h.cartService.DeleteShoppingCartItem(service.DeleteShoppingCartItemDTO{
		CorrelationID:      correlationID,
		ShoppingCartItemID: req.ShoppingCartItemID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item removed from wishlist"})
}
