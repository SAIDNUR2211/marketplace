package db

import (
	"marketplace/internal/models/domain"
	"time"
)

var (
	UserRole      = "USER"
	AdminRole     = "ADMIN"
	ShopkeperRole = "SHOPKEPER"
)

type User struct {
	ID        int       `db:"id"`
	Username  string    `db:"username"`
	FullName  string    `db:"full_name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Role      string    `db:"role"`
	Phone     string    `db:"phone"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt time.Time `db:"deleted_at"`
}

func (u *User) ToDomain() domain.User {
	return domain.User{
		ID:        u.ID,
		FullName:  u.FullName,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		Role:      u.Role,
		Phone:     u.Phone,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,
	}
}
func (u *User) FromDomain(d *domain.User) {
	u.ID = d.ID
	u.FullName = d.FullName
	u.Username = d.Username
	u.Email = d.Email
	u.Password = d.Password
	u.Role = d.Role
	u.Phone = d.Phone
	u.CreatedAt = d.CreatedAt
	u.UpdatedAt = d.UpdatedAt
	u.DeletedAt = d.DeletedAt
}
