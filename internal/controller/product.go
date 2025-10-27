package controller

import (
	"marketplace/internal/errs"
	"marketplace/internal/models/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateProductHandler godoc
// @Summary Создать продукт
// @Description Создает новый продукт (только для админов и владельцев магазинов)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body domain.Product true "Данные продукта"
// @Success 201 {object} domain.Product
// @Failure 400 {object} CommonError
// @Failure 401 {object} CommonError
// @Failure 403 {object} CommonError
// @Failure 422 {object} CommonError
// @Router /api/v1/products [post]
func (ctrl *Controller) CreateProductHandler(c *gin.Context) {
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errs.ErrInvalidRequestBody.Error()})
		return
	}
	err := ctrl.service.CreateProduct(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, product)
}

// GetProductByIDHandler godoc
// @Summary Получить продукт по ID
// @Description Получает информацию о продукте по его ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} domain.Product
// @Failure 400 {object} CommonError
// @Failure 404 {object} CommonError
// @Router /api/v1/products/{id} [get]
func (ctrl *Controller) GetProductByIDHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": errs.ErrInvalidProductID.Error()})
		return
	}
	product, err := ctrl.service.GetProductByID(id)
	if err != nil {
		if err == errs.ErrProductNotfound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

// UpdateProductHandler godoc
// @Summary Обновить продукт
// @Description Обновляет информацию о продукте (только для админов и владельцев магазинов)
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param input body domain.Product true "Данные для обновления"
// @Success 200 {object} domain.Product
// @Failure 400 {object} CommonError
// @Failure 401 {object} CommonError
// @Failure 403 {object} CommonError
// @Failure 404 {object} CommonError
// @Router /api/v1/products/{id} [put]
func (ctrl *Controller) UpdateProductHandler(c *gin.Context) {
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errs.ErrInvalidRequestBody.Error()})
		return
	}
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": errs.ErrInvalidProductID.Error()})
		return
	}
	product.ID = id
	userIDUntyped, _ := c.Get(userIDCtx)
	userRoleUntyped, _ := c.Get(userRoleCtx)
	userID := userIDUntyped.(int)
	userRole := userRoleUntyped.(string)

	product.ID = id
	if err := ctrl.service.UpdateProduct(&product, userID, userRole); err != nil {
		ctrl.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, product)
}

// DeleteProductHandler godoc
// @Summary Удалить продукт
// @Description Удаляет продукт (soft delete, только для админов и владельцев магазинов)
// @Tags products
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 200 {object} CommonResponse
// @Failure 400 {object} CommonError
// @Failure 401 {object} CommonError
// @Failure 403 {object} CommonError
// @Failure 404 {object} CommonError
// @Router /api/v1/products/{id} [delete]
func (ctrl *Controller) DeleteProductHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": errs.ErrInvalidProductID.Error()})
		return
	}

	userIDUntyped, _ := c.Get(userIDCtx)
	userRoleUntyped, _ := c.Get(userRoleCtx)
	userID := userIDUntyped.(int)
	userRole := userRoleUntyped.(string)

	if err := ctrl.service.DeleteProduct(id, userID, userRole); err != nil {
		ctrl.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}

// ListProductsHandler godoc
// @Summary Список продуктов
// @Description Получает список продуктов магазина
// @Tags products
// @Produce json
// @Param shop_id query int true "Shop ID"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} domain.Product
// @Failure 400 {object} CommonError
// @Failure 500 {object} CommonError
// @Router /api/v1/products [get]
func (ctrl *Controller) ListProductsHandler(c *gin.Context) {
	shopIDParam := c.Query("shop_id")
	limitParam := c.DefaultQuery("limit", "20")
	offsetParam := c.DefaultQuery("offset", "0")
	shopID, err := strconv.ParseInt(shopIDParam, 10, 64)
	if err != nil || shopID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": errs.ErrInvalidFieldValue.Error()})
		return
	}
	limit, _ := strconv.Atoi(limitParam)
	offset, _ := strconv.Atoi(offsetParam)
	products, err := ctrl.service.ListProducts(shopID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}
