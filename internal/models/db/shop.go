package db

import (
	"marketplace/internal/models/domain"
	"time"
)

type Shop struct {
	ID          int64      `db:"id"`
	Name        string     `db:"name"`
	Slug        string     `db:"slug"`
	OwnerID     int64      `db:"owner_id"`
	Description *string    `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"` // ИСПРАВЛЕНО: указатель
}

func (sh *Shop) ToDomain() *domain.Shop {
	domainShop := &domain.Shop{
		ID:        sh.ID,
		Name:      sh.Name,
		Slug:      sh.Slug,
		OwnerID:   sh.OwnerID,
		CreatedAt: sh.CreatedAt,
		UpdatedAt: sh.UpdatedAt,
	}

	if sh.Description != nil {
		domainShop.Description = *sh.Description
	}
	if sh.DeletedAt != nil {
		domainShop.DeletedAt = *sh.DeletedAt
	}

	return domainShop
}

func (sh *Shop) FromDomain(d *domain.Shop) {
	sh.ID = d.ID
	sh.Name = d.Name
	sh.Slug = d.Slug
	sh.OwnerID = d.OwnerID
	sh.Description = &d.Description
	sh.CreatedAt = d.CreatedAt
	sh.UpdatedAt = d.UpdatedAt

	if !d.DeletedAt.IsZero() {
		sh.DeletedAt = &d.DeletedAt
	}
}
