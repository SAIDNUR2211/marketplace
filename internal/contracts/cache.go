package contracts

import (
	"context"
	"marketplace/internal/models/domain"
	"time"
)

type CacheI interface {
	// Product methods
	GetProduct(ctx context.Context, key string) (*domain.Product, error)
	SetProduct(ctx context.Context, key string, product *domain.Product, expiration time.Duration) error
	DelProduct(ctx context.Context, key string) error

	// Shop methods
	GetShop(ctx context.Context, key string) (*domain.Shop, error)
	SetShop(ctx context.Context, key string, shop *domain.Shop, expiration time.Duration) error
	DelShop(ctx context.Context, key string) error

	// User methods
	GetUser(ctx context.Context, key string) (*domain.User, error)
	SetUser(ctx context.Context, key string, user *domain.User, expiration time.Duration) error
	DelUser(ctx context.Context, key string) error

	// Common methods
	Ping(ctx context.Context) error
	Close() error
}
