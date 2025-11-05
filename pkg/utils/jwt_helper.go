package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GetClaimsFromContext(c echo.Context) (*JWTClaims, error) {
	user := c.Get("user")
	if user == nil {
		return nil, errors.New("missing user in context")
	}

	token, ok := user.(*jwt.Token)
	if !ok {
		return nil, errors.New("invalid token type")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("invalid claims type")
	}

	return claims, nil
}
