package main

import (
	"giat-cerika-service/configs"
	"giat-cerika-service/routes"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	db := configs.InitDB()
	defer configs.CloseDB(db)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	routes.InitRoutes(e, db)

	// Start server
	e.Logger.Fatal(e.Start(":" + configs.GetConfig("PORT")))
}
