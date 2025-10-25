package routes

import (
	"giat-cerika-service/internal/handlers"
	"giat-cerika-service/internal/repositories"
	"giat-cerika-service/internal/services"
	"giat-cerika-service/pkg/utils"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func InitRoutes(e *echo.Echo, db *gorm.DB) {
	// Validator
	e.Validator = utils.NewValidator()

	// Repositories
	adminRepo := repositories.NewAdminRepository(db)
	studentRepo := repositories.NewStudentRepository(db)

	// Services
	adminService := services.NewAdminService(adminRepo)
	studentService := services.NewStudentService(studentRepo)

	// Handlers
	adminHandler := handlers.NewAdminHandler(adminService)
	studentHandler := handlers.NewStudentHandler(studentService)

	// API Group
	api := e.Group("/api/v1")

	// Admin Routes
	InitAdminRoutes(api, adminHandler)

	// Student Routes
	InitStudentRoutes(api, studentHandler)
}
