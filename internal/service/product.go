package service

import (
	"context" // НОВЫЙ ИМПОРТ
	"errors"
	"fmt" // НОВЫЙ ИМПОРТ
	"marketplace/internal/errs"
	"marketplace/internal/models/domain"
	"marketplace/utils"
	"time" // НОВЫЙ ИМПОРТ

	"github.com/go-redis/redis/v8" // НОВЫЙ ИМПОРТ
)

func (s *Service) CreateProduct(product *domain.Product) error {
	s.logger.Info().Str("func", "CreateProduct").Msg("creating product")
	if product.Price <= 0 {
		return errs.ErrInvalidFieldValue
	}
	if product.ShopID == 0 {
		return errs.ErrInvalidFieldValue
	}

	product.Slug = utils.GenerateSlug(product.Name)

	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	if err := s.repository.CreateProduct(product); err != nil {
		s.logger.Error().Err(err).Msg("failed to create product")
		return err
	}
	s.logger.Info().Int64("id", product.ID).Msg("product created successfully")
	go func() {
		cacheKey := fmt.Sprintf("product:%d", product.ID)
		s.logger.Info().Int64("id", product.ID).Msg("caching: setting NEW product to redis")
		if err := s.cache.SetProduct(context.Background(), cacheKey, product, time.Hour*1); err != nil {
			s.logger.Error().Err(err).Int64("id", product.ID).Msg("caching new product failed")
		}
	}()
	return nil
}

func (s *Service) GetProductByID(id int64) (*domain.Product, error) {
	s.logger.Info().Int64("id", id).Msg("fetching product by id")
	if id <= 0 {
		return nil, errs.ErrInvalidProductID
	}

	// --- НОВАЯ ЛОГИКА КЕШИРОВАНИЯ ---
	ctx := context.Background()
	// 1. Формируем ключ для Redis
	cacheKey := fmt.Sprintf("product:%d", id)

	// 2. Пытаемся достать из Redis
	cachedProduct, err := s.cache.GetProduct(ctx, cacheKey)
	if err == nil {
		// 2a. Cache Hit (Нашли в кеше)
		s.logger.Info().Int64("id", id).Msg("cache hit: product found in redis")
		return cachedProduct, nil
	}

	// 2b. Если ошибка - это НЕ "не найдено", значит Redis упал.
	if !errors.Is(err, redis.Nil) {
		s.logger.Error().Err(err).Int64("id", id).Msg("redis error on get product")
		// (Мы не останавливаемся, а просто идем в БД, как будто кеша нет)
	} else {
		// 2c. Cache Miss (Не нашли в кеше)
		s.logger.Info().Int64("id", id).Msg("cache miss: product not found in redis")
	}

	// 3. Идем в базу данных (PostgreSQL)
	product, err := s.repository.GetProductByID(id)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get product by id from db")
		if err == errs.ErrNotfound {
			return nil, errs.ErrProductNotfound
		}
		return nil, err
	}

	// 4. Сохраняем в кеш ДЛЯ СЛЕДУЮЩЕГО РАЗА.
	// Мы делаем это в горутине, чтобы пользователь не ждал ответа от Redis.
	go func() {
		s.logger.Info().Int64("id", id).Msg("caching: setting product to redis")
		err := s.cache.SetProduct(context.Background(), cacheKey, product, time.Hour*1) // Кеш на 1 час
		if err != nil {
			s.logger.Error().Err(err).Int64("id", id).Msg("caching failed: could not set product to redis")
		}
	}()
	// --- КОНЕЦ ЛОГИКИ КЕШИРОВАНИЯ ---

	return product, nil
}

func (s *Service) UpdateProduct(product *domain.Product, userID int, userRole string) error {
	s.logger.Info().Int64("id", product.ID).Msg("updating product")
	if product == nil || product.ID <= 0 {
		return errs.ErrInvalidProductID
	}

	existingProduct, err := s.repository.GetProductByID(product.ID)
	if err != nil {
		return err
	}
	shop, err := s.repository.GetShopByID(existingProduct.ShopID)
	if err != nil {
		return err
	}

	if userRole != domain.AdminRole && shop.OwnerID != int64(userID) {
		return errors.New("permission denied: you are not the owner of this product's shop")
	}

	if product.Name != "" && len(product.Name) < 4 {
		return errs.ErrInvalidProductName
	}
	if product.Price < 0 {
		return errs.ErrInvalidFieldValue
	}

	product.UpdatedAt = time.Now()
	if err := s.repository.UpdateProduct(product); err != nil {
		s.logger.Error().Err(err).Msg("failed to update product")
		return err
	}

	// --- НОВЫЙ КОД: ИНВАЛИДАЦИЯ КЕША ---
	// Мы обновили товар, значит старая запись в кеше неверна.
	// Удаляем ее. В горутине, чтобы не блокировать.
	go func() {
		cacheKey := fmt.Sprintf("product:%d", product.ID)
		s.logger.Info().Int64("id", product.ID).Msg("cache invalidation: deleting product from redis")
		if err := s.cache.DelProduct(context.Background(), cacheKey); err != nil {
			s.logger.Error().Err(err).Int64("id", product.ID).Msg("cache invalidation failed")
		}
	}()
	// --- КОНЕЦ ---

	s.logger.Info().Int64("id", product.ID).Msg("product updated successfully")
	return nil
}

func (s *Service) DeleteProduct(id int64, userID int, userRole string) error {
	s.logger.Info().Int64("id", id).Msg("deleting product")
	if id <= 0 {
		return errs.ErrInvalidProductID
	}

	existingProduct, err := s.repository.GetProductByID(id)
	if err != nil {
		return err
	}
	shop, err := s.repository.GetShopByID(existingProduct.ShopID)
	if err != nil {
		return err
	}

	if userRole != domain.AdminRole && shop.OwnerID != int64(userID) {
		return errors.New("permission denied: you are not the owner of this product's shop")
	}

	if err := s.repository.DeleteProduct(id); err != nil {
		s.logger.Error().Err(err).Msg("failed to delete product")
		return err
	}

	// --- НОВЫЙ КОД: ИНВАЛИДАЦИЯ КЕША ---
	// Товар удален (soft delete), удаляем его из кеша.
	go func() {
		cacheKey := fmt.Sprintf("product:%d", id)
		s.logger.Info().Int64("id", id).Msg("cache invalidation: deleting product from redis")
		if err := s.cache.DelProduct(context.Background(), cacheKey); err != nil {
			s.logger.Error().Err(err).Int64("id", id).Msg("cache invalidation failed")
		}
	}()
	// --- КОНЕЦ ---

	s.logger.Info().Int64("id", id).Msg("product soft deleted successfully")
	return nil
}

func (s *Service) ListProducts(shopID int64, limit, offset int) ([]*domain.Product, error) {
	s.logger.Info().Int64("shop_id", shopID).Int("limit", limit).Int("offset", offset).Msg("listing products")
	if shopID <= 0 {
		return nil, errs.ErrInvalidFieldValue
	}
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	products, err := s.repository.ListProducts(shopID, limit, offset)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to list products")
		return nil, err
	}
	s.logger.Info().Int("count", len(products)).Msg("products listed successfully")
	return products, nil
}
