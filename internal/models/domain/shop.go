package domain

import "time"

// Shop represents a shop
// @Description Shop information
type Shop struct {
	ID          int64     `json:"id" example:"1"`
	Name        string    `json:"name" example:"My Shop"`
	Slug        string    `json:"slug" example:"my-shop"`
	OwnerID     int64     `json:"owner_id" example:"123"`
	Description string    `json:"description" example:"Shop description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}
