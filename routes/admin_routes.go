package routes

import (
	"giat-cerika-service/internal/handlers"
	"giat-cerika-service/internal/middlewares"

	"github.com/labstack/echo/v4"
)

func InitAdminRoutes(g *echo.Group, h *handlers.AdminHandler) {
	admin := g.Group("/admin")

	// Public routes
	admin.POST("/login", h.Login)

	// Protected routes
	admin.GET("/profile", h.GetProfile, middlewares.JWTMiddleware, middlewares.AdminOnly)
	admin.PUT("/profile", h.UpdateProfile, middlewares.JWTMiddleware, middlewares.AdminOnly)
	admin.POST("/logout", h.Logout, middlewares.JWTMiddleware, middlewares.AdminOnly)
}
