package db

import (
	"marketplace/internal/models/domain"
	"time"
)

type Order struct {
	ID        int64      `db:"id"`
	UserID    int64      `db:"user_id"`
	ShopID    *int64     `db:"shop_id"`
	Total     float64    `db:"total"`
	Currency  string     `db:"currency"`
	Status    string     `db:"status"`
	Note      *string    `db:"note"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (o *Order) ToDomain() *domain.Order {
	domainOrder := &domain.Order{
		ID:        o.ID,
		UserID:    o.UserID,
		ShopID:    o.ShopID,
		Total:     o.Total,
		Currency:  o.Currency,
		Status:    o.Status,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
		DeletedAt: o.DeletedAt,
	}
	if o.Note != nil {
		domainOrder.Note = *o.Note
	}
	return domainOrder
}

func (o *Order) FromDomain(d *domain.Order) {
	o.ID = d.ID
	o.UserID = d.UserID
	o.ShopID = d.ShopID
	o.Total = d.Total
	o.Status = d.Status
	o.Note = &d.Note
	o.CreatedAt = d.CreatedAt
	o.UpdatedAt = d.UpdatedAt
	o.DeletedAt = d.DeletedAt
}
