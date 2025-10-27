package domain

import "time"

// Product represents a product
// @Description Product information
type Product struct {
	ID          int64      `json:"id" example:"1"`
	SKU         *string    `json:"sku,omitempty" example:"SKU123"`
	Name        string     `json:"name" example:"Product Name"`
	Slug        string     `json:"slug" example:"product-name"`
	Description *string    `json:"description,omitempty" example:"Product description"`
	Price       float64    `json:"price" example:"99.99"`
	Currency    string     `json:"currency" example:"USD"`
	Quantity    int        `json:"quantity" example:"10"`
	ShopID      int64      `json:"shop_id" example:"1"`
	Active      bool       `json:"active" example:"true"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
