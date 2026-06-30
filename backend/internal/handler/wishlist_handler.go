package handler

import (
	"bobs-used-bookstore-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WishlistHandler struct {
	cartService *service.ShoppingCartService
}

func NewWishlistHandler(cartService *service.ShoppingCartService) *WishlistHandler {
	return &WishlistHandler{cartService: cartService}
}

func (h *WishlistHandler) Index(c *gin.Context) {
	correlationID := c.GetString("cartCorrelationID")
	cart, err := h.cartService.GetOrCreateCart(correlationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var wishlistItems []interface{}
	for _, item := range cart.Items {
		if !item.WantToBuy {
			wishlistItems = append(wishlistItems, item)
		}
	}

	c.JSON(http.StatusOK, gin.H{"items": wishlistItems})
}

func (h *WishlistHandler) MoveToCart(c *gin.Context) {
	itemID, err := strconv.Atoi(c.PostForm("shoppingCartItemId"))
	if err != nil {
		var body struct {
			ShoppingCartItemID int `json:"shoppingCartItemId"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
			return
		}
		itemID = body.ShoppingCartItemID
	}

	correlationID := c.GetString("cartCorrelationID")
	cart, err := h.cartService.GetOrCreateCart(correlationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.cartService.MoveToCart(itemID, cart.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "moved to cart"})
}

func (h *WishlistHandler) MoveAllToCart(c *gin.Context) {
	correlationID := c.GetString("cartCorrelationID")
	cart, err := h.cartService.GetOrCreateCart(correlationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.cartService.MoveAllToCart(cart.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "all items moved to cart"})
}

func (h *WishlistHandler) Delete(c *gin.Context) {
	itemID, err := strconv.Atoi(c.PostForm("shoppingCartItemId"))
	if err != nil {
		var body struct {
			ShoppingCartItemID int `json:"shoppingCartItemId"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
			return
		}
		itemID = body.ShoppingCartItemID
	}

	correlationID := c.GetString("cartCorrelationID")
	cart, err := h.cartService.GetOrCreateCart(correlationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.cartService.RemoveItem(itemID, cart.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item removed from wishlist"})
}
