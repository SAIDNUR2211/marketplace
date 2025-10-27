package controller

import (
	"marketplace/internal/errs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SetRoleRequest struct {
	Role string `json:"role" binding:"required" example:"SHOPKEPER"`
}

// SetUserRoleHandler godoc
// @Summary Изменить роль пользователя
// @Description Изменяет роль пользователя (только для админов)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param input body SetRoleRequest true "Новая роль"
// @Success 200 {object} CommonResponse
// @Failure 400 {object} CommonError
// @Failure 401 {object} CommonError
// @Failure 403 {object} CommonError
// @Failure 404 {object} CommonError
// @Router /api/v1/admin/users/{id}/role [put]
func (ctrl *Controller) SetUserRoleHandler(c *gin.Context) {
	actorUserID, _ := c.Get(userIDCtx)
	actorRole, _ := c.Get(userRoleCtx)

	targetUserID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ctrl.handleError(c, errs.ErrInvalidID)
		return
	}

	var req SetRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.handleError(c, errs.ErrInvalidRequestBody)
		return
	}

	err = ctrl.service.SetUserRole(actorUserID.(int), actorRole.(string), targetUserID, req.Role)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, CommonResponse{Message: "User role updated successfully"})
}
