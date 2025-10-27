package repository

import (
	"marketplace/internal/errs"
	"marketplace/internal/models/db"
	"marketplace/internal/models/domain"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

func (repository *Repository) CreateUser(user domain.User) (err error) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("func_name", "repository.CreateUser").Logger()

	_, err = repository.db.Exec(
		`INSERT INTO users (full_name, username, email, password, phone) VALUES ($1, $2, $3, $4, $5)`,
		user.FullName, user.Username, user.Email, user.Password, user.Phone,
	)

	if err != nil {
		logger.Err(err).Msg("error inserting user")
		// ИСПРАВЛЕНИЕ: добавляем проверку на нарушение уникальности
		if strings.Contains(err.Error(), "unique constraint") {
			if strings.Contains(err.Error(), "username") {
				return errs.ErrUsernameAlreadyExists
			} else if strings.Contains(err.Error(), "email") {
				return errs.ErrEmailAlreadyExists
			} else if strings.Contains(err.Error(), "phone") {
				return errs.ErrPhoneAlreadyExists
			}
		}
		return repository.translateError(err)
	}
	return nil
}
func (repository *Repository) GetUserByID(id int) (domain.User, error) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("func_name", "repository.GetUserByID").Logger()
	var dbUser db.User
	if err := repository.db.Get(&dbUser, `SELECT id, full_name, username, email, password, role,  phone, created_at, updated_at FROM users WHERE id = $1`, id); err != nil {
		logger.Err(err).Msg("error selecting user")
		return domain.User{}, repository.translateError(err)
	}
	return dbUser.ToDomain(), nil
}
func (repository *Repository) GetUserByUsername(username string) (domain.User, error) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("func_name", "repository.GetUserByUsername").Logger()
	var dbUser db.User
	if err := repository.db.Get(&dbUser, `SELECT id, full_name, username, email, password, role,  phone, created_at, updated_at FROM users WHERE username = $1 AND deleted_at IS NULL`, username); err != nil {
		logger.Err(err).Msg("error selecting user")
		return domain.User{}, repository.translateError(err)
	}
	return dbUser.ToDomain(), nil
}
func (repository *Repository) GetUserByEmail(email string) (domain.User, error) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("func_name", "repository.GetUserByEmail").Logger()
	var dbUser db.User
	if err := repository.db.Get(&dbUser, `SELECT id, full_name, username, email, password, role,  phone, created_at, updated_at FROM users WHERE email = $1`, email); err != nil {
		logger.Err(err).Msg("error selecting user")
		return domain.User{}, repository.translateError(err)
	}
	return dbUser.ToDomain(), nil
}
func (repository *Repository) GetUserByRole(role string) (domain.User, error) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("func_name", "repository.GetUserByEmail").Logger()
	var dbUser db.User
	if err := repository.db.Get(&dbUser, `SELECT id, full_name, username, email, password, role,  phone, created_at, updated_at FROM users WHERE role = $1`, role); err != nil {
		logger.Err(err).Msg("error selecting user")
		return domain.User{}, repository.translateError(err)
	}
	return dbUser.ToDomain(), nil
}
func (repository *Repository) GetUserByPhone(phone string) (domain.User, error) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("func_name", "repository.GetUserByEmail").Logger()
	var dbUser db.User
	if err := repository.db.Get(&dbUser, `SELECT id, full_name, username, email, password, role,  phone, created_at, updated_at FROM users WHERE phone = $1 AND deleted_at IS NULL`, phone); err != nil {
		logger.Err(err).Msg("error selecting user")
		return domain.User{}, repository.translateError(err)
	}
	return dbUser.ToDomain(), nil
}
func (repository *Repository) UpdateUserRole(userID int, role string) error {
	logger := zerolog.New(os.Stdout).With().Timestamp().Str("func_name", "repository.UpdateUserRole").Logger()
	query := `UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2`
	result, err := repository.db.Exec(query, role, userID)
	if err != nil {
		logger.Err(err).Msg("error updating user role")
		return repository.translateError(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return repository.translateError(err)
	}
	if rowsAffected == 0 {
		return errs.ErrUserNotFound
	}
	return nil
}
