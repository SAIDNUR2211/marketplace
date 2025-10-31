package contracts

import (
	"marketplace/internal/models/domain"

	"github.com/jmoiron/sqlx"
)

type RepositoryI interface {
	CreateUser(user domain.User) (err error)
	GetUserByID(id int) (domain.User, error)
	GetUserByUsername(username string) (domain.User, error)
	GetUserByEmail(email string) (domain.User, error)
	GetUserByRole(role string) (domain.User, error)
	GetUserByPhone(phone string) (domain.User, error)
	UpdateUserRole(userID int, role string) error
	CreateProduct(product *domain.Product) error
	GetProductByID(id int64) (*domain.Product, error)
	UpdateProduct(product *domain.Product) error
	DeleteProduct(id int64) error
	ListProducts(shopID int64, limit, offset int) ([]*domain.Product, error)
	DecreaseProductQuantity(productID int64, quantity int) error
	CreateShop(shop *domain.Shop) error
	GetShopByID(id int64) (*domain.Shop, error)
	UpdateShop(shop *domain.Shop) error
	DeleteShop(id int64) error
	ListShops(ownerID int64, limit, offset int) ([]*domain.Shop, error)
	GetOrderByID(orderID int64) (*domain.Order, []domain.OrderItem, error)
	BeginTx() (*sqlx.Tx, error)
	CreateOrderWithTx(tx *sqlx.Tx, order *domain.Order, items []domain.OrderItem) (int64, error)
	GetProductByIDWithTx(tx *sqlx.Tx, id int64) (*domain.Product, error)
	DecreaseProductQuantityWithTx(tx *sqlx.Tx, productID int64, quantity int) error
}
