package service

import (
	"errors"
	"marketplace/internal/errs"
	"marketplace/internal/models/domain"
	"marketplace/utils"
)

func (s *Service) CreateUser(user domain.User) (err error) {
	user.Role = domain.UserRole

	_, err = s.repository.GetUserByUsername(user.Username)
	if err == nil {
		return errs.ErrUsernameAlreadyExists
	} else if !errors.Is(err, errs.ErrNotfound) {
		return err
	}

	_, err = s.repository.GetUserByEmail(user.Email)
	if err == nil {
		return errs.ErrEmailAlreadyExists
	} else if !errors.Is(err, errs.ErrNotfound) {
		return err
	}

	_, err = s.repository.GetUserByPhone(user.Phone)
	if err == nil {
		return errs.ErrPhoneAlreadyExists
	} else if !errors.Is(err, errs.ErrNotfound) {
		return err
	}

	user.Password, err = utils.GenerateHash(user.Password)
	if err != nil {
		return err
	}

	if err = s.repository.CreateUser(user); err != nil {
		return err
	}
	return nil
}

func (s *Service) Authenticate(user domain.User) (int, string, error) {
	if user.Username == "" && user.Email == "" {
		return 0, "", errs.ErrIncorrectUsernameOrPassword
	}
	var userFromDB domain.User
	var err error
	if user.Username != "" && user.Email != "" {
		return 0, "", errors.New("use either username or email for login, not both")
	}
	if user.Username != "" {
		userFromDB, err = s.repository.GetUserByUsername(user.Username)
		if err != nil {
			if errors.Is(err, errs.ErrNotfound) {
				return 0, "", errs.ErrIncorrectUsernameOrPassword
			}
			return 0, "", err
		}
	} else {
		userFromDB, err = s.repository.GetUserByEmail(user.Email)
		if err != nil {
			if errors.Is(err, errs.ErrNotfound) {
				return 0, "", errs.ErrIncorrectUsernameOrPassword
			}
			return 0, "", err
		}
	}
	if !utils.CheckPasswordHash(user.Password, userFromDB.Password) {
		return 0, "", errs.ErrIncorrectUsernameOrPassword
	}
	return userFromDB.ID, userFromDB.Role, nil
}
func (s *Service) SetUserRole(actorUserID int, actorRole string, targetUserID int, newRole string) error {
	if actorRole != domain.AdminRole {
		return errors.New("permission denied: only admins can change user roles")
	}
	validRoles := map[string]bool{
		domain.UserRole:      true,
		domain.ShopkeperRole: true,
		domain.AdminRole:     true,
	}

	if !validRoles[newRole] {
		return errs.ErrInvalidFieldValue
	}
	if actorUserID == targetUserID && newRole != domain.AdminRole {
		return errors.New("admins cannot remove their own admin role")
	}
	_, err := s.repository.GetUserByID(targetUserID)
	if err != nil {
		if errors.Is(err, errs.ErrNotfound) {
			return errs.ErrUserNotFound
		}
		return err
	}
	return s.repository.UpdateUserRole(targetUserID, newRole)
}
