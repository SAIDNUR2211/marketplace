package domain

import "time"

// Order represents an order
// @Description Order information
type Order struct {
	ID        int64       `json:"id"`
	UserID    int64       `json:"user_id"`
	ShopID    *int64      `json:"shop_id,omitempty"`
	Total     float64     `json:"total"`
	Currency  string      `json:"currency"`
	Status    string      `json:"status"`
	Note      string      `json:"note,omitempty"`
	Items     []OrderItem `json:"items,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	DeletedAt *time.Time  `json:"deleted_at,omitempty"`
}

// CreateOrderInput represents input for creating an order
// @Description Input for creating an order
type CreateOrderInput struct {
	Note  string                 `json:"note"`
	Items []CreateOrderItemInput `json:"items"`
}

// CreateOrderItemInput represents input for order item
// @Description Input for order item
type CreateOrderItemInput struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}
