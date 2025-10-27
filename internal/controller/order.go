package controller

import (
	"errors"
	"marketplace/internal/errs"
	"marketplace/internal/models/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateOrderHandler godoc
// @Summary      Создать заказ
// @Description  Создает новый заказ для аутентифицированного пользователя
// @Tags         orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input body domain.CreateOrderInput true "Данные заказа"
// @Success      201  {object}  map[string]int64
// @Failure      400  {object}  CommonError
// @Failure      401  {object}  CommonError
// @Failure      500  {object}  CommonError
// @Router       /api/v1/orders [post]
func (ctrl *Controller) CreateOrderHandler(c *gin.Context) {
	userIDUntyped, exists := c.Get(userIDCtx)
	if !exists {
		ctrl.handleError(c, errs.ErrInvalidToken)
		return
	}
	userID := userIDUntyped.(int)

	var input domain.CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		ctrl.handleError(c, errs.ErrInvalidRequestBody)
		return
	}

	orderID, err := ctrl.service.CreateOrder(userID, input)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"order_id": orderID})
}

// GetOrderHandler godoc
// @Summary      Получить заказ по ID
// @Description  Получает детали заказа по его ID
// @Tags         orders
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Order ID"
// @Success      200  {object}  domain.Order
// @Failure      400  {object}  CommonError
// @Failure      401  {object}  CommonError
// @Failure      403  {object}  CommonError
// @Failure      404  {object}  CommonError
// @Router       /api/v1/orders/{id} [get]
func (ctrl *Controller) GetOrderHandler(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		ctrl.handleError(c, errs.ErrInvalidID)
		return
	}

	order, items, err := ctrl.service.GetOrderByID(orderID)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}

	userIDUntyped, _ := c.Get(userIDCtx)
	userRole := c.GetString(userRoleCtx)
	userID := userIDUntyped.(int)

	if userRole != domain.AdminRole && order.UserID != int64(userID) {
		ctrl.handleError(c, errors.New("permission denied"))
		return
	}

	order.Items = items

	c.JSON(http.StatusOK, order)
}
