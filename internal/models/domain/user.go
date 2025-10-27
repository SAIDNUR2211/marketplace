package domain

import "time"

var (
	UserRole      = "USER"
	AdminRole     = "ADMIN"
	ShopkeperRole = "SHOPKEPER"
)

// User represents a user
// @Description User information

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
