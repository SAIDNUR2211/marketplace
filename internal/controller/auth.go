package controller

import (
	"errors"
	"marketplace/internal/errs"
	"marketplace/internal/models/domain"
	"marketplace/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SignUpRequest struct {
	FullName string `json:"full_name" example:"John Doe"`
	Username string `json:"username" example:"johndoe"`
	Password string `json:"password" example:"password123"`
	Email    string `json:"email" example:"john@example.com"`
	Phone    string `json:"phone" example:"+1234567890"`
}

// SignUp godoc
// @Summary      Регистрация пользователя
// @Description  Создание нового пользователя
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body SignUpRequest true "Данные для регистрации"
// @Success      201  {object}  CommonResponse
// @Failure      400  {object}  CommonError
// @Failure      422  {object}  CommonError
// @Failure      500  {object}  CommonError
// @Router       /auth/sign-up [post]
func (ctrl *Controller) SignUp(c *gin.Context) {
	var input SignUpRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		ctrl.handleError(c, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}
	if input.Email == "" || input.Password == "" || input.Username == "" {
		ctrl.handleError(c, errors.New("email, password and username are required"))
		return
	}
	if err := ctrl.service.CreateUser(domain.User{
		FullName: input.FullName,
		Username: input.Username,
		Password: input.Password,
		Email:    input.Email,
		Phone:    input.Phone,
	}); err != nil {
		ctrl.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, CommonResponse{Message: "User created successfully!"})
}

type SignInRequest struct {
	Username string `json:"username" example:"johndoe"`
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"password123"`
}

type TokenPairResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// SignIn godoc
// @Summary      Вход в систему
// @Description  Аутентификация пользователя и получение токенов
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body SignInRequest true "Данные для входа"
// @Success      200  {object}  TokenPairResponse
// @Failure      400  {object}  CommonError
// @Failure      401  {object}  CommonError
// @Router       /auth/sign-in [post]
func (ctrl *Controller) SignIn(c *gin.Context) {
	var input SignInRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		ctrl.handleError(c, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	// Валидация - должно быть указано только одно поле
	if (input.Username == "" && input.Email == "") || (input.Username != "" && input.Email != "") {
		ctrl.handleError(c, errors.New("please provide either username or email for login"))
		return
	}

	if input.Password == "" {
		ctrl.handleError(c, errors.New("password is required"))
		return
	}

	userID, userRole, err := ctrl.service.Authenticate(domain.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		ctrl.handleError(c, err)
		return
	}

	accessToken, refreshToken, err := ctrl.generateNewTokenPair(userID, userRole)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, TokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

const (
	refreshTokenHeader = "X-Refresh-Token"
)

// RefreshTokenPair godoc
// @Summary      Обновление токенов
// @Description  Получение новой пары access/refresh токенов
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        X-Refresh-Token header string true "Refresh Token" default(Bearer {token})
// @Success      200  {object}  TokenPairResponse
// @Failure      401  {object}  CommonError
// @Router       /auth/refresh [get]
func (ctrl *Controller) RefreshTokenPair(c *gin.Context) {
	token, err := ctrl.extractTokenFromHeader(c, refreshTokenHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, CommonError{Error: err.Error()})
		return
	}
	userID, isRefresh, userRole, err := pkg.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, CommonError{Error: err.Error()})
		return
	}
	if !isRefresh {
		c.JSON(http.StatusUnauthorized, CommonError{Error: "inappropriate token"})
		return
	}
	accessToken, refreshToken, err := ctrl.generateNewTokenPair(userID, userRole)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, TokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
