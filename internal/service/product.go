package service

import (
	"errors"
	"marketplace/internal/errs"
	"marketplace/internal/models/domain"
	"marketplace/utils"
	"time"
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
	return nil
}
func (s *Service) GetProductByID(id int64) (*domain.Product, error) {
	s.logger.Info().Int64("id", id).Msg("fetching product by id")
	if id <= 0 {
		return nil, errs.ErrInvalidProductID
	}
	product, err := s.repository.GetProductByID(id)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get product by id from db")
		if err == errs.ErrNotfound {
			return nil, errs.ErrProductNotfound
		}
		return nil, err
	}
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
