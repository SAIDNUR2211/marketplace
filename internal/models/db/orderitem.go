package db

import (
	"marketplace/internal/models/domain"
	"time"
)

type OrderItem struct {
	ID         int64     `db:"id"`
	OrderID    int64     `db:"order_id"`
	ProductID  int64     `db:"product_id"`
	Name       string    `db:"name"`
	SKU        string    `db:"sku"`
	UnitPrice  float64   `db:"unit_price"`
	Quantity   int       `db:"quantity"`
	TotalPrice float64   `db:"total_price"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

func (o *OrderItem) ToDomain() *domain.OrderItem {
	return &domain.OrderItem{
		ID:         o.ID,
		OrderID:    o.OrderID,
		ProductID:  o.ProductID,
		Name:       o.Name,
		SKU:        o.SKU,
		UnitPrice:  o.UnitPrice,
		Quantity:   o.Quantity,
		TotalPrice: o.TotalPrice,
		CreatedAt:  o.CreatedAt,
		UpdatedAt:  o.UpdatedAt,
	}
}
func (o *OrderItem) FromDomain(d *domain.OrderItem) {
	o.ID = d.ID
	o.OrderID = d.OrderID
	o.ProductID = d.ProductID
	o.Name = d.Name
	o.SKU = d.SKU
	o.UnitPrice = d.UnitPrice
	o.Quantity = d.Quantity
	o.TotalPrice = d.TotalPrice
	o.CreatedAt = d.CreatedAt
	o.UpdatedAt = d.UpdatedAt
}
