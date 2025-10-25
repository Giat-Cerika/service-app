package routes

import (
	"giat-cerika-service/internal/handlers"
	"giat-cerika-service/internal/middlewares"

	"github.com/labstack/echo/v4"
)

func InitStudentRoutes(g *echo.Group, h *handlers.StudentHandler) {
	student := g.Group("/student")

	// Public routes
	student.POST("/register", h.Register)
	student.POST("/login", h.Login)

	// Protected routes
	student.GET("/profile", h.GetProfile, middlewares.JWTMiddleware)
	student.PUT("/profile", h.UpdateProfile, middlewares.JWTMiddleware)
	student.POST("/logout", h.Logout, middlewares.JWTMiddleware)
}
