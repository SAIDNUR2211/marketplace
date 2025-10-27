package service

import (
	"errors"
	"fmt"
	"marketplace/internal/errs"
	"marketplace/internal/models/domain"
)

func (s *Service) CreateOrder(userID int, input domain.CreateOrderInput) (int64, error) {
	if len(input.Items) == 0 {
		return 0, errors.New("order must contain at least one item")
	}
	tx, err := s.repository.BeginTx()
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to begin transaction")
		return 0, err
	}
	var committed bool
	defer func() {
		if !committed {
			if rbErr := tx.Rollback(); rbErr != nil {
				s.logger.Error().Err(rbErr).Msg("failed to rollback transaction")
			}
		}
	}()
	var total float64
	var currency string
	orderItems := make([]domain.OrderItem, 0, len(input.Items))
	productMap := make(map[int64]*domain.Product)
	for _, item := range input.Items {
		product, err := s.repository.GetProductByIDWithTx(tx, item.ProductID)
		if err != nil {
			if errors.Is(err, errs.ErrNotfound) {
				return 0, errs.ErrProductNotfound
			}
			s.logger.Error().Err(err).Int64("product_id", item.ProductID).Msg("failed to get product with tx")
			return 0, err
		}
		if product.Quantity < item.Quantity {
			return 0, fmt.Errorf("not enough stock for product: %s (available: %d, requested: %d)",
				product.Name, product.Quantity, item.Quantity)
		}
		if currency == "" {
			currency = product.Currency
		} else if currency != product.Currency {
			return 0, errors.New("all items in order must have the same currency")
		}
		itemTotal := product.Price * float64(item.Quantity)
		total += itemTotal
		orderItems = append(orderItems, domain.OrderItem{
			ProductID:  product.ID,
			Name:       product.Name,
			UnitPrice:  product.Price,
			Quantity:   item.Quantity,
			TotalPrice: itemTotal,
		})
		productMap[product.ID] = product
	}
	newOrder := &domain.Order{
		UserID:   int64(userID),
		Total:    total,
		Currency: currency,
		Status:   "pending",
		Note:     input.Note,
	}
	orderID, err := s.repository.CreateOrderWithTx(tx, newOrder, orderItems)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create order with tx")
		return 0, err
	}
	for _, item := range orderItems {
		err := s.repository.DecreaseProductQuantityWithTx(tx, item.ProductID, item.Quantity)
		if err != nil {
			s.logger.Error().Err(err).Int64("product_id", item.ProductID).Msg("failed to decrease product quantity with tx")
			return 0, fmt.Errorf("failed to update stock for product %d: %w", item.ProductID, err)
		}
	}
	if err := tx.Commit(); err != nil {
		s.logger.Error().Err(err).Msg("failed to commit transaction")
		return 0, err
	}
	committed = true
	s.logger.Info().Int64("order_id", orderID).Msg("order created successfully")
	return orderID, nil
}
func (s *Service) GetOrderByID(orderID int64) (*domain.Order, []domain.OrderItem, error) {
	return s.repository.GetOrderByID(orderID)
}
