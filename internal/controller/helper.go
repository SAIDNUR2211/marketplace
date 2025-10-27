package controller

import (
	"errors"
	"marketplace/internal/configs"
	"marketplace/pkg"
	"strings"

	"github.com/gin-gonic/gin"
)

func (ctrl *Controller) extractTokenFromHeader(c *gin.Context, headerKey string) (string, error) {
	header := c.GetHeader(headerKey)

	if header == "" {
		return "", errors.New("empty authorization header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid authorization header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("empty token")
	}

	return headerParts[1], nil
}

func (ctrl *Controller) generateNewTokenPair(userID int, userRole string) (string, string, error) {
	accessToken, err := pkg.GenerateToken(userID,
		configs.AppSettings.AuthParams.AccessTokenTtlMinutes,
		userRole, false)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := pkg.GenerateToken(userID,
		configs.AppSettings.AuthParams.RefreshTokenTtlDays,
		userRole, true)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
