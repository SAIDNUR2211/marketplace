package domain

import "time"

// OrderItem represents an item in an order
// @Description Order item information
type OrderItem struct {
	ID         int64     `json:"id"`
	ProductID  int64     `json:"product_id"`
	OrderID    int64     `json:"order_id"`
	Name       string    `json:"name,omitempty"`
	SKU        string    `json:"sku,omitempty"`
	UnitPrice  float64   `json:"unit_price"`
	Quantity   int       `json:"quantity"`
	TotalPrice float64   `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
