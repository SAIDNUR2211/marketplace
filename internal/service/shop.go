package service

import (
	"errors"
	"marketplace/internal/errs"
	"marketplace/internal/models/domain"
	"marketplace/utils"
	"time"
)

func (s *Service) CreateShop(shop *domain.Shop) error {
	if shop == nil {
		return errs.ErrInvalidRequestBody
	}
	if shop.OwnerID <= 0 {
		return errs.ErrInvalidFieldValue
	}

	shop.Slug = utils.GenerateSlug(shop.Name)

	shop.CreatedAt = time.Now()
	shop.UpdatedAt = time.Now()
	if err := s.repository.CreateShop(shop); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetShopByID(id int64) (*domain.Shop, error) {
	if id <= 0 {
		return nil, errs.ErrInvalidShopID
	}
	shop, err := s.repository.GetShopByID(id)
	if err != nil {
		if err == errs.ErrNotfound {
			return nil, errs.ErrShopNotFound
		}
		return nil, err
	}
	return shop, nil
}
func (s *Service) UpdateShop(shop *domain.Shop, userID int, userRole string) error {
	if shop.ID <= 0 {
		return errs.ErrInvalidShopID
	}

	existingShop, err := s.repository.GetShopByID(shop.ID)
	if err != nil {
		return err
	}

	if userRole != domain.AdminRole && existingShop.OwnerID != int64(userID) {
		return errors.New("permission denied: you are not the owner of this shop")
	}

	shop.UpdatedAt = time.Now()
	if err := s.repository.UpdateShop(shop); err != nil {
		return err
	}
	return nil
}
func (s *Service) DeleteShop(id int64, userID int, userRole string) error {
	if id <= 0 {
		return errs.ErrInvalidShopID
	}

	existingShop, err := s.repository.GetShopByID(id)
	if err != nil {
		return err
	}

	if userRole != domain.AdminRole && existingShop.OwnerID != int64(userID) {
		return errors.New("permission denied: you are not the owner of this shop")
	}

	if err := s.repository.DeleteShop(id); err != nil {
		return err
	}
	return nil
}
func (s *Service) ListShops(ownerID int64, limit, offset int) ([]*domain.Shop, error) {
	if ownerID <= 0 {
		return nil, errs.ErrInvalidFieldValue
	}
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	shops, err := s.repository.ListShops(ownerID, limit, offset)
	if err != nil {
		return nil, err
	}
	return shops, nil
}
