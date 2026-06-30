package handler

import (
	"bobs-used-bookstore-api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ShoppingCartHandler struct {
	cartService *service.ShoppingCartService
}

func NewShoppingCartHandler(cartService *service.ShoppingCartService) *ShoppingCartHandler {
	return &ShoppingCartHandler{cartService: cartService}
}

func (h *ShoppingCartHandler) Index(c *gin.Context) {
	correlationID := c.GetString("cartCorrelationID")
	cart, err := h.cartService.GetOrCreateCart(correlationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var cartItems []interface{}
	var subTotal float64
	for _, item := range cart.Items {
		if item.WantToBuy {
			cartItems = append(cartItems, item)
			if item.Book != nil && item.Book.IsInStock() {
				subTotal += item.Book.Price
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"items":    cartItems,
		"subTotal": subTotal,
	})
}

func (h *ShoppingCartHandler) Delete(c *gin.Context) {
	itemID, err := strconv.Atoi(c.PostForm("shoppingCartItemId"))
	if err != nil {
		// Try JSON body
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

	c.JSON(http.StatusOK, gin.H{"message": "item removed"})
}
