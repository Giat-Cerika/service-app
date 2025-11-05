package middlewares

import (
	"giat-cerika-service/pkg/utils"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func RoleMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user")
			if user == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing user in context")
			}

			token, ok := user.(*jwt.Token)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token type")
			}

			claims, ok := token.Claims.(*utils.JWTClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid claims type")
			}

			for _, role := range allowedRoles {
				if claims.Role == role {
					return next(c)
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "forbidden: insufficient role permission")
		}
	}
}
