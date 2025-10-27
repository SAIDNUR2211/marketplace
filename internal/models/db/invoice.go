package db

import "time"

type Invoice struct {
	ID        int64      `db:"id"`
	OrderID   int64      `db:"order_id"`
	Amount    float64    `db:"amount"`
	Currency  string     `db:"currency"`
	Paid      bool       `db:"paid"`
	Method    *string    `db:"method"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	PaidAt    *time.Time `db:"paid_at"`
}
