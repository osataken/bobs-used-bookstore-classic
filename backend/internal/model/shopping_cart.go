package model

type ShoppingCart struct {
	Entity
	CorrelationID string             `gorm:"type:text" json:"correlationId"`
	Items         []ShoppingCartItem `gorm:"foreignKey:ShoppingCartID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}
