package middlewares

import (
	"context"
	"fmt"
	"giat-cerika-service/configs"
	"giat-cerika-service/pkg/utils"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	jwtmiddleware "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func JWTMiddleware(rdb *redis.Client) echo.MiddlewareFunc {
	return jwtmiddleware.WithConfig(jwtmiddleware.Config{
		SigningKey: []byte(configs.GetJWTSecret()),
		ContextKey: "user",

		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(utils.JWTClaims)
		},

		Skipper: func(c echo.Context) bool {
			token := c.Request().Header.Get("Authorization")
			if token != "" {
				token = strings.TrimPrefix(token, "Bearer ")
				ctx := context.Background()
				val, err := rdb.Get(ctx, fmt.Sprintf("blacklist:%s", token)).Result()
				if err == nil && val == "blacklisted" {
					return true
				}
			}
			return false
		},

		ErrorHandler: func(c echo.Context, err error) error {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		},
	})
}

func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}
