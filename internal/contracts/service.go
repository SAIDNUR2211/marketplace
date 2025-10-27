package contracts

import "marketplace/internal/models/domain"

type ServiceI interface {
	CreateUser(user domain.User) error
	Authenticate(user domain.User) (int, string, error)
	CreateProduct(product *domain.Product) error
	SetUserRole(actorUserID int, actorRole string, targetUserID int, newRole string) error
	GetProductByID(id int64) (*domain.Product, error)
	UpdateProduct(product *domain.Product, userID int, userRole string) error
	DeleteProduct(id int64, userID int, userRole string) error
	ListProducts(shopID int64, limit, offset int) ([]*domain.Product, error)
	CreateShop(shop *domain.Shop) error
	GetShopByID(id int64) (*domain.Shop, error)
	UpdateShop(shop *domain.Shop, userID int, userRole string) error
	DeleteShop(id int64, userID int, userRole string) error
	ListShops(ownerID int64, limit, offset int) ([]*domain.Shop, error)
	CreateOrder(userID int, input domain.CreateOrderInput) (int64, error)
	GetOrderByID(orderID int64) (*domain.Order, []domain.OrderItem, error)
}
