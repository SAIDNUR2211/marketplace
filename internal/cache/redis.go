package cache

import (
	"context"
	"encoding/json"
	"marketplace/internal/configs"
	"marketplace/internal/contracts"
	"marketplace/internal/models/domain"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
}

func InitRedis(params configs.RedisParams) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     params.Addr,
		Password: params.Password,
		DB:       params.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewCache(client *redis.Client) contracts.CacheI {
	return &RedisCache{
		client: client,
	}
}

func (r *RedisCache) GetProduct(ctx context.Context, key string) (*domain.Product, error) {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, redis.Nil
		}
		return nil, err
	}

	var product domain.Product
	err = json.Unmarshal([]byte(data), &product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *RedisCache) SetProduct(ctx context.Context, key string, product *domain.Product, expiration time.Duration) error {
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, expiration).Err()
}

func (r *RedisCache) DelProduct(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) GetShop(ctx context.Context, key string) (*domain.Shop, error) {
	return nil, redis.Nil
}

func (r *RedisCache) SetShop(ctx context.Context, key string, shop *domain.Shop, expiration time.Duration) error {
	return nil
}

func (r *RedisCache) DelShop(ctx context.Context, key string) error {
	return nil
}

func (r *RedisCache) GetUser(ctx context.Context, key string) (*domain.User, error) {
	return nil, redis.Nil
}

func (r *RedisCache) SetUser(ctx context.Context, key string, user *domain.User, expiration time.Duration) error {
	return nil
}

func (r *RedisCache) DelUser(ctx context.Context, key string) error {
	return nil
}

func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}
