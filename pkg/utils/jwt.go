package utils

import (
	"giat-cerika-service/configs"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GenerateToken(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(configs.GetConfig("JWT_SECRET")))
}

func GetUserIDFromToken(c echo.Context) uint {
	return c.Get("user_id").(uint)
}

func GetRoleFromToken(c echo.Context) string {
	return c.Get("role").(string)
}
