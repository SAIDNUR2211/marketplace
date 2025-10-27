package controller

import (
	"marketplace/internal/errs"
	"marketplace/internal/models/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateShopHandler godoc
// @Summary Создать магазин
// @Description Создает новый магазин
// @Tags shops
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body domain.Shop true "Данные магазина"
// @Success 201 {object} domain.Shop
// @Failure 400 {object} CommonError
// @Failure 401 {object} CommonError
// @Failure 422 {object} CommonError
// @Router /api/v1/shops [post]
func (ctrl *Controller) CreateShopHandler(c *gin.Context) {
	var shop domain.Shop
	if err := c.ShouldBindJSON(&shop); err != nil {
		ctrl.handleError(c, errs.ErrInvalidRequestBody)
		return
	}
	if err := ctrl.service.CreateShop(&shop); err != nil {
		ctrl.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, shop)
}

// GetShopByIDHandler godoc
// @Summary Получить магазин по ID
// @Description Получает информацию о магазине по его ID
// @Tags shops
// @Produce json
// @Param id path int true "Shop ID"
// @Success 200 {object} domain.Shop
// @Failure 400 {object} CommonError
// @Failure 404 {object} CommonError
// @Router /api/v1/shops/{id} [get]
func (ctrl *Controller) GetShopByIDHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		ctrl.handleError(c, errs.ErrInvalidShopID)
		return
	}
	shop, err := ctrl.service.GetShopByID(id)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, shop)
}

// UpdateShopHandler godoc
// @Summary Обновить магазин
// @Description Обновляет информацию о магазине (только для админов и владельцев)
// @Tags shops
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Shop ID"
// @Param input body domain.Shop true "Данные для обновления"
// @Success 200 {object} domain.Shop
// @Failure 400 {object} CommonError
// @Failure 401 {object} CommonError
// @Failure 403 {object} CommonError
// @Failure 404 {object} CommonError
// @Router /api/v1/shops/{id} [put]
func (ctrl *Controller) UpdateShopHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		ctrl.handleError(c, errs.ErrInvalidShopID)
		return
	}

	userIDUntyped, _ := c.Get(userIDCtx)
	userRoleUntyped, _ := c.Get(userRoleCtx)
	userID := userIDUntyped.(int)
	userRole := userRoleUntyped.(string)

	var shop domain.Shop
	if err := c.ShouldBindJSON(&shop); err != nil {
		ctrl.handleError(c, errs.ErrInvalidRequestBody)
		return
	}
	shop.ID = id

	if err := ctrl.service.UpdateShop(&shop, userID, userRole); err != nil {
		ctrl.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, shop)
}

// DeleteShopHandler godoc
// @Summary Удалить магазин
// @Description Удаляет магазин (soft delete, только для админов и владельцев)
// @Tags shops
// @Produce json
// @Security BearerAuth
// @Param id path int true "Shop ID"
// @Success 200 {object} CommonResponse
// @Failure 400 {object} CommonError
// @Failure 401 {object} CommonError
// @Failure 403 {object} CommonError
// @Failure 404 {object} CommonError
// @Router /api/v1/shops/{id} [delete]
func (ctrl *Controller) DeleteShopHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		ctrl.handleError(c, errs.ErrInvalidShopID)
		return
	}

	userIDUntyped, _ := c.Get(userIDCtx)
	userRoleUntyped, _ := c.Get(userRoleCtx)
	userID := userIDUntyped.(int)
	userRole := userRoleUntyped.(string)

	if err := ctrl.service.DeleteShop(id, userID, userRole); err != nil {
		ctrl.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "shop deleted successfully"})
}

// ListShopsHandler godoc
// @Summary Список магазинов
// @Description Получает список магазинов владельца
// @Tags shops
// @Produce json
// @Param owner_id query int true "Owner ID"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} domain.Shop
// @Failure 400 {object} CommonError
// @Failure 500 {object} CommonError
// @Router /api/v1/shops [get]
func (ctrl *Controller) ListShopsHandler(c *gin.Context) {
	ownerIDParam := c.Query("owner_id")
	ownerID, err := strconv.ParseInt(ownerIDParam, 10, 64)
	if err != nil || ownerID <= 0 {
		ctrl.handleError(c, errs.ErrInvalidFieldValue)
		return
	}
	limit := 20
	offset := 0
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}
	shops, err := ctrl.service.ListShops(ownerID, limit, offset)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, shops)
}
