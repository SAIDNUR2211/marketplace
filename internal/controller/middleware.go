package controller

import (
	"marketplace/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userIDCtx           = "userID"
	userRoleCtx         = "userRole"
)

func (ctrl *Controller) checkUserAuthentication(c *gin.Context) {
	token, err := ctrl.extractTokenFromHeader(c, authorizationHeader)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, CommonError{Error: err.Error()})
		return
	}
	userID, isRefresh, userRole, err := pkg.ParseToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, CommonError{Error: err.Error()})
		return
	}
	if isRefresh {
		c.AbortWithStatusJSON(http.StatusUnauthorized, CommonError{Error: "inappropriate token"})
		return
	}
	c.Set(userIDCtx, userID)
	c.Set(userRoleCtx, string(userRole))

}

func (ctrl *Controller) checkRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get(userRoleCtx)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, CommonError{Error: "user role not found in context"})
			return
		}

		for _, role := range allowedRoles {
			if role == userRole.(string) {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, CommonError{Error: "permission denied for this role"})
	}
}
