package main

import (
	"giat-cerika-service/configs"
	datasources "giat-cerika-service/internal/dataSources"
	"giat-cerika-service/internal/middlewares"
	"giat-cerika-service/pkg/workers/producer"
	"giat-cerika-service/routes"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	configs.LoadEnv()

	db := configs.InitDB()
	rdb := configs.InitRedis()

	configs.RunMigrations(db)

	e := echo.New()
	e.Use(middlewares.LoggerMiddleware)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
	}))

	for _, r := range e.Routes() {
		log.Printf("ROUTE %s %s", r.Method, r.Path)
	}

	cloudinarySvc, err := datasources.NewCloudinaryService()
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary service: %v", err)
	}

	configs.InitRabbitMQ()
	defer configs.CloseConnections()

	go producer.StartWorker()

	routes.Routes(e, db, rdb, &cloudinarySvc)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(e.Start(":" + port))
}
