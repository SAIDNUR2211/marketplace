package db

import (
	"marketplace/internal/models/domain"
	"time"
)

type Product struct {
	ID          int64      `db:"id"`
	SKU         *string    `db:"sku"`
	Name        string     `db:"name"`
	Slug        string     `db:"slug"`
	Description *string    `db:"description"`
	Price       float64    `db:"price"`
	Currency    string     `db:"currency"`
	Quantity    int        `db:"quantity"`
	ShopID      int64      `db:"shop_id"`
	Active      bool       `db:"active"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

func (p *Product) ToDomain() *domain.Product {
	return &domain.Product{
		ID:          p.ID,
		SKU:         p.SKU,
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		Price:       p.Price,
		Currency:    p.Currency,
		Quantity:    p.Quantity,
		ShopID:      p.ShopID,
		Active:      p.Active,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		DeletedAt:   p.DeletedAt,
	}
}
func (p *Product) FromDomain(d *domain.Product) {
	p.ID = d.ID
	p.SKU = d.SKU
	p.Name = d.Name
	p.Slug = d.Slug
	p.Description = d.Description
	p.Price = d.Price
	p.Currency = d.Currency
	p.Quantity = d.Quantity
	p.ShopID = d.ShopID
	p.Active = d.Active
	p.CreatedAt = d.CreatedAt
	p.UpdatedAt = d.UpdatedAt
	p.DeletedAt = d.DeletedAt
}
