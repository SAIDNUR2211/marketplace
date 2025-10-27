package repository

import (
	"database/sql"
	"marketplace/internal/models/db"
	"marketplace/internal/models/domain"

	"github.com/jmoiron/sqlx"
)

func (r *Repository) CreateOrderWithTx(tx *sqlx.Tx, order *domain.Order, items []domain.OrderItem) (int64, error) {
	var orderID int64
	orderQuery := `INSERT INTO orders (user_id, total, currency, status, note) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := tx.Get(&orderID, orderQuery, order.UserID, order.Total, order.Currency, order.Status, order.Note)
	if err != nil {
		return 0, r.translateError(err)
	}
	itemQuery := `INSERT INTO order_items (order_id, product_id, name, unit_price, quantity, total_price) VALUES ($1, $2, $3, $4, $5, $6)`
	for _, item := range items {
		_, err = tx.Exec(itemQuery, orderID, item.ProductID, item.Name, item.UnitPrice, item.Quantity, item.TotalPrice)
		if err != nil {
			return 0, r.translateError(err)
		}
	}
	return orderID, nil
}
func (r *Repository) GetOrderByID(orderID int64) (*domain.Order, []domain.OrderItem, error) {
	var dbOrder db.Order
	queryOrder := `SELECT id, user_id, total, currency, status, note, created_at, updated_at FROM orders WHERE id=$1 AND deleted_at IS NULL`
	if err := r.db.Get(&dbOrder, queryOrder, orderID); err != nil {
		return nil, nil, r.translateError(err)
	}
	var dbItems []db.OrderItem
	queryItems := `SELECT id, order_id, product_id, name, unit_price, quantity, total_price, created_at, updated_at FROM order_items WHERE order_id=$1`
	if err := r.db.Select(&dbItems, queryItems, orderID); err != nil {
		if err == sql.ErrNoRows {
			return dbOrder.ToDomain(), []domain.OrderItem{}, nil
		}
		return nil, nil, r.translateError(err)
	}
	domainItems := make([]domain.OrderItem, len(dbItems))
	for i, item := range dbItems {
		domainItems[i] = *item.ToDomain()
	}
	return dbOrder.ToDomain(), domainItems, nil
}
