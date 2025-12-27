package studentroute

import (
	datasources "giat-cerika-service/internal/dataSources"
	studenthandler "giat-cerika-service/internal/handlers/student_handler"
	"giat-cerika-service/internal/middlewares"
	classrepo "giat-cerika-service/internal/repositories/class_repo"
	studentrepo "giat-cerika-service/internal/repositories/student_repo"
	studentservice "giat-cerika-service/internal/services/student_service"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func StudentRoutes(e *echo.Group, db *gorm.DB, rdb *redis.Client, cld *datasources.CloudinaryService) {
	studentRepo := studentrepo.NewStudentRepositoryImpl(db)
	classRepo := classrepo.NewClassRepositoryImpl(db)
	studentService := studentservice.NewStudentServiceImpl(studentRepo, classRepo, rdb, *cld)
	studentHandler := studenthandler.NewStudentHandler(studentService)

	e.POST("/register", studentHandler.RegisterStudent)
	e.POST("/login", studentHandler.LoginStudent)
	e.POST("/check-nisn-and-dateofbirth", studentHandler.CheckNisnAndDateOfBirthStudent)
	e.PUT("/update-new-password", studentHandler.UpdateNewPasswordStudent)

	studentGroup := e.Group("", middlewares.JWTMiddleware(rdb), middlewares.RoleMiddleware("student"))
	studentGroup.GET("/me", studentHandler.GetProfileStudent, middlewares.JWTMiddleware(rdb))
	studentGroup.POST("/logout", studentHandler.Logout, middlewares.JWTMiddleware(rdb))
	studentGroup.PUT("/update-profile", studentHandler.UpdateProfileStudent)
	studentGroup.PUT("/edit-photo", studentHandler.EditPhotoStudent)
	studentGroup.POST("/tooth-brush", studentHandler.CreateToothBrush)
	studentGroup.GET("/history-tooth-brush", studentHandler.GetHistoryToothBrush)
	studentGroup.GET("/all", studentHandler.GetStudentAll)
}
