package domain

import "time"

type Invoice struct {
	ID        int64     `json:"id"`
	OrderID   int64     `json:"order_id"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	Paid      bool      `json:"paid"`
	Method    string    `json:"method,omitempty"`
	PaidAt    string    `json:"paid_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
