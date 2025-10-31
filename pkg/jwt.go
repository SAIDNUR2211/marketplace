package pkg

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	jwt.StandardClaims
	UserID    int    `json:"user_id"`
	Role      string `json:"role"`
	IsRefresh bool   `json:"is_refresh"`
}

func GenerateToken(userID, ttl int, Role string, isRefresh bool) (string, error) {
	claims := CustomClaims{
		StandardClaims: jwt.StandardClaims{},
		UserID:         userID,
		IsRefresh:      isRefresh,
		Role:           Role,
	}
	if isRefresh {
		claims.StandardClaims.ExpiresAt = time.Now().Add(time.Duration(ttl) * 24 * time.Hour).Unix()
	} else {
		claims.StandardClaims.ExpiresAt = time.Now().Add(time.Duration(ttl) * time.Minute).Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
func ParseToken(tokenString string) (int, bool, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return 0, false, "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.UserID, claims.IsRefresh, claims.Role, nil
	}
	return 0, false, "", fmt.Errorf("invalid token")
}
