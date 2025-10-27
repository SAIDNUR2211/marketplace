package repository

import (
	"marketplace/internal/errs"
	"marketplace/internal/models/db"
	"marketplace/internal/models/domain"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func (r *Repository) CreateShop(shop *domain.Shop) error {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("func", "CreateShop").Logger()
	dbShop := db.Shop{}
	dbShop.FromDomain(shop)
	now := time.Now()
	query := `INSERT INTO shops (name, slug, owner_id, description, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(query, dbShop.Name, dbShop.Slug, dbShop.OwnerID, dbShop.Description, now, now).Scan(&dbShop.ID, &dbShop.CreatedAt, &dbShop.UpdatedAt)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create shop")
		return errs.ErrSomethingWentWrong
	}
	shop.ID = dbShop.ID
	shop.CreatedAt = dbShop.CreatedAt
	shop.UpdatedAt = dbShop.UpdatedAt
	logger.Info().Int64("shop_id", shop.ID).Msg("shop created successfully")
	return nil
}
func (r *Repository) GetShopByID(id int64) (*domain.Shop, error) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("func", "GetShopByID").Logger()
	if id <= 0 {
		return nil, errs.ErrInvalidID
	}
	var dbShop db.Shop
	query := `SELECT id, name, slug, owner_id, description, created_at, updated_at, deleted_at FROM shops WHERE id = $1 AND deleted_at IS NULL`
	err := r.db.Get(&dbShop, query, id)
	if err != nil {
		logger.Error().Err(err).Int64("id", id).Msg("failed to get shop by id")
		return nil, errs.ErrNotfound
	}
	return dbShop.ToDomain(), nil
}
func (r *Repository) UpdateShop(shop *domain.Shop) error {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("func", "UpdateShop").Logger()
	if shop == nil || shop.ID <= 0 {
		return errs.ErrInvalidID
	}
	dbShop := db.Shop{}
	dbShop.FromDomain(shop)
	dbShop.UpdatedAt = time.Now()
	query := `UPDATE shops SET name = $1, slug = $2, description = $3, updated_at = $4 WHERE id = $5 AND deleted_at IS NULL`
	_, err := r.db.Exec(query, dbShop.Name, dbShop.Slug, dbShop.Description, dbShop.UpdatedAt, dbShop.ID)
	if err != nil {
		logger.Error().Err(err).Int64("id", dbShop.ID).Msg("failed to update shop")
		return errs.ErrSomethingWentWrong
	}
	logger.Info().Int64("shop_id", dbShop.ID).Msg("shop updated successfully")
	return nil
}
func (r *Repository) DeleteShop(id int64) error {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("func", "DeleteShop").Logger()
	if id <= 0 {
		return errs.ErrInvalidID
	}
	query := `UPDATE shops SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		logger.Error().Err(err).Int64("id", id).Msg("failed to delete shop")
		return errs.ErrSomethingWentWrong
	}
	logger.Info().Int64("shop_id", id).Msg("shop soft deleted successfully")
	return nil
}
func (r *Repository) ListShops(ownerID int64, limit, offset int) ([]*domain.Shop, error) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("func", "ListShops").Logger()
	if ownerID <= 0 {
		return nil, errs.ErrInvalidID
	}
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	var dbShops []db.Shop
	// ИСПРАВЛЕНИЕ: правильная параметризация
	query := `SELECT id, name, slug, owner_id, description, created_at, updated_at, deleted_at 
              FROM shops 
              WHERE owner_id = $1 AND deleted_at IS NULL 
              ORDER BY created_at DESC 
              LIMIT $2 OFFSET $3`

	err := r.db.Select(&dbShops, query, ownerID, limit, offset)
	if err != nil {
		logger.Error().Err(err).Msg("failed to list shops")
		return nil, r.translateError(err)
	}

	shops := make([]*domain.Shop, 0, len(dbShops))
	for _, sh := range dbShops {
		shops = append(shops, sh.ToDomain())
	}
	logger.Info().Int("count", len(shops)).Msg("shops listed successfully")
	return shops, nil
}
