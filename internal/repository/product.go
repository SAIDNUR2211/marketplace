package repository

import (
	"errors"
	"marketplace/internal/models/db"
	"marketplace/internal/models/domain"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

func (r *Repository) CreateProduct(product *domain.Product) error {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("func", "CreateProduct").Logger()
	dbProduct := db.Product{}
	dbProduct.FromDomain(product)
	query := `INSERT INTO products (sku, name, slug, description, price, currency, quantity, shop_id, active, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id, created_at, updated_at`
	now := time.Now()
	err := r.db.QueryRow(
		query, dbProduct.SKU, dbProduct.Name, dbProduct.Slug, dbProduct.Description, dbProduct.Price, dbProduct.Currency, dbProduct.Quantity, dbProduct.ShopID, dbProduct.Active, now, now).Scan(&dbProduct.ID, &dbProduct.CreatedAt, &dbProduct.UpdatedAt)

	if err != nil {
		logger.Error().Err(err).Msg("failed to create product")
		return r.translateError(err)
	}
	product.ID = dbProduct.ID
	product.CreatedAt = dbProduct.CreatedAt
	product.UpdatedAt = dbProduct.UpdatedAt
	logger.Info().Int64("product_id", product.ID).Msg("product created successfully")
	return nil
}

func (r *Repository) GetProductByID(id int64) (*domain.Product, error) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("func", "GetProductByID").Logger()
	var dbProduct db.Product
	query := `SELECT id, sku, name, slug, description, price, currency, quantity, shop_id, active, created_at, updated_at, deleted_at FROM products WHERE id = $1 AND deleted_at IS NULL`
	err := r.db.Get(&dbProduct, query, id)
	if err != nil {
		logger.Error().Err(err).Int64("id", id).Msg("failed to get product by id")
		return nil, r.translateError(err)
	}
	return dbProduct.ToDomain(), nil
}
func (r *Repository) UpdateProduct(product *domain.Product) error {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("func", "UpdateProduct").Logger()
	dbProduct := db.Product{}
	dbProduct.FromDomain(product)
	query := `UPDATE products SET sku = $1, name = $2, slug = $3, description = $4, price = $5, currency = $6, quantity = $7, shop_id = $8, active = $9, updated_at = $10 WHERE id = $11 AND deleted_at IS NULL`
	_, err := r.db.Exec(
		query,
		dbProduct.SKU,
		dbProduct.Name,
		dbProduct.Slug,
		dbProduct.Description,
		dbProduct.Price,
		dbProduct.Currency,
		dbProduct.Quantity,
		dbProduct.ShopID,
		dbProduct.Active,
		time.Now(),
		dbProduct.ID,
	)
	if err != nil {
		logger.Error().Err(err).Int64("id", dbProduct.ID).Msg("failed to update product")
		return r.translateError(err)
	}
	logger.Info().Int64("product_id", dbProduct.ID).Msg("product updated successfully")
	return nil
}
func (r *Repository) DeleteProduct(id int64) error {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("func", "DeleteProduct").Logger()
	query := `UPDATE products SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		logger.Error().Err(err).Int64("id", id).Msg("failed to delete product")
		return r.translateError(err)
	}
	logger.Info().Int64("product_id", id).Msg("product soft deleted successfully")
	return nil
}
func (r *Repository) ListProducts(shopID int64, limit, offset int) ([]*domain.Product, error) {
	logger := zerolog.New(os.Stdout).With().
		Timestamp().Str("func", "ListProducts").Logger()
	var dbProducts []db.Product

	// ИСПРАВЛЕНИЕ: правильная параметризация
	query := `SELECT id, sku, name, slug, description, price, currency, quantity, shop_id, active, created_at, updated_at, deleted_at 
              FROM products 
              WHERE shop_id = $1 AND deleted_at IS NULL 
              ORDER BY created_at DESC 
              LIMIT $2 OFFSET $3`

	err := r.db.Select(&dbProducts, query, shopID, limit, offset)
	if err != nil {
		logger.Error().Err(err).Msg("failed to list products")
		return nil, r.translateError(err)
	}

	products := make([]*domain.Product, 0, len(dbProducts))
	for _, p := range dbProducts {
		products = append(products, p.ToDomain())
	}
	logger.Info().Int("count", len(products)).Msg("products listed successfully")
	return products, nil
}
func (r *Repository) DecreaseProductQuantity(productID int64, quantity int) error {
	query := `UPDATE products SET quantity = quantity - $1 WHERE id = $2 AND quantity >= $1`

	result, err := r.db.Exec(query, quantity, productID)
	if err != nil {
		return r.translateError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return r.translateError(err)
	}

	if rowsAffected == 0 {
		return errors.New("not enough stock")
	}

	return nil
}
func (r *Repository) GetProductByIDWithTx(tx *sqlx.Tx, id int64) (*domain.Product, error) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("func", "GetProductByIDWithTx").Logger()
	var dbProduct db.Product

	// ИЗМЕНЕНИЕ: Добавлен 'FOR UPDATE' для блокировки строки
	query := `SELECT id, sku, name, slug, description, price, currency, quantity, shop_id, active, created_at, updated_at, deleted_at 
	          FROM products WHERE id = $1 AND deleted_at IS NULL FOR UPDATE`

	// ИЗМЕНЕНИЕ: Используем tx.Get()
	err := tx.Get(&dbProduct, query, id)
	if err != nil {
		logger.Error().Err(err).Int64("id", id).Msg("failed to get product by id with tx")
		return nil, r.translateError(err)
	}
	return dbProduct.ToDomain(), nil
}

// DecreaseProductQuantityWithTx уменьшает кол-во товара внутри транзакции.
func (r *Repository) DecreaseProductQuantityWithTx(tx *sqlx.Tx, productID int64, quantity int) error {
	query := `UPDATE products SET quantity = quantity - $1 WHERE id = $2 AND quantity >= $1`

	// ИЗМЕНЕНИЕ: Используем tx.Exec()
	result, err := tx.Exec(query, quantity, productID)
	if err != nil {
		return r.translateError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return r.translateError(err)
	}

	if rowsAffected == 0 {
		// Эта ошибка будет поймана сервисом и вызовет Rollback
		return errors.New("not enough stock")
	}

	return nil
}
