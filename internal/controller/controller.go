package controller

import (
	"errors"
	"marketplace/internal/contracts"
	"marketplace/internal/errs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service contracts.ServiceI
}

func NewController(svc contracts.ServiceI) *Controller {
	return &Controller{
		service: svc,
	}
}

// handleError godoc
// @Description Обрабатывает ошибки и возвращает соответствующий HTTP статус
func (ctrl *Controller) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, errs.ErrProductNotfound) ||
		errors.Is(err, errs.ErrUserNotFound) ||
		errors.Is(err, errs.ErrNotfound):
		c.JSON(http.StatusNotFound, CommonError{Error: err.Error()})
	case errors.Is(err, errs.ErrInvalidProductID) || errors.Is(err, errs.ErrInvalidRequestBody):
		c.JSON(http.StatusBadRequest, CommonError{Error: err.Error()})
	case errors.Is(err, errs.ErrIncorrectUsernameOrPassword) || errors.Is(err, errs.ErrInvalidToken):
		c.JSON(http.StatusUnauthorized, CommonError{Error: err.Error()})
	case errors.Is(err, errs.ErrInvalidFieldValue) ||
		errors.Is(err, errs.ErrInvalidProductName) ||
		errors.Is(err, errs.ErrUsernameAlreadyExists):
		c.JSON(http.StatusUnprocessableEntity, CommonError{Error: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, CommonError{Error: err.Error()})
	}
}

// ping godoc
// @Summary Ping
// @Description Проверка работоспособности API
// @Tags health
// @Produce json
// @Success 200 {object} CommonResponse
// @Router /ping [get]
func (ctrl *Controller) ping(c *gin.Context) {
	c.JSON(http.StatusOK, CommonResponse{
		Message: "PONG",
	})
}

// healthCheck godoc
// @Summary Health Check
// @Description Проверка статуса сервиса
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (ctrl *Controller) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "SHOPE",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
