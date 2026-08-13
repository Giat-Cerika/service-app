package main

import (
	"giat-cerika-service/configs"
	datasources "giat-cerika-service/internal/dataSources"
	"giat-cerika-service/pkg/workers/producer"
	"giat-cerika-service/routes"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	configs.LoadEnv()

	db := configs.InitDB()
	rdb := configs.InitRedis()

	configs.RunMigrations(db)

	e := echo.New()

	// ================================================================
	// MIDDLEWARE STACK — urutan penting!
	// ================================================================

	// 1. Recover — tangkap panic agar server tidak crash
	e.Use(middleware.Recover())

	// 2. Request Timeout — 30 detik maks per request
	//    Mencegah koneksi menggantung dan menghabiskan resource
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout:      30 * time.Second,
		ErrorMessage: "request timeout, silakan coba lagi",
	}))

	// 3. Body Limit — maksimal 25MB per request (untuk menampung foto kamera HP)
	e.Use(middleware.BodyLimit("25M"))

	// 4. Gzip — kompresi response untuk hemat bandwidth
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	// 5. Rate Limiter — mencegah abuse/DDoS
	//    20 request/detik per IP (cukup untuk 300 user)
	e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      20,
				Burst:     40,
				ExpiresIn: 3 * time.Minute,
			},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			return ctx.RealIP(), nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, map[string]string{
				"message": "terlalu banyak request, coba lagi nanti",
			})
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, map[string]string{
				"message": "terlalu banyak request, coba lagi nanti",
			})
		},
	}))

	// 6. CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*", "https://website-giat-cerika.vercel.app"},
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.DELETE,
			echo.PATCH,
			echo.OPTIONS,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
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

	// ================================================================
	// HTTP SERVER — tuning untuk high concurrency
	// ================================================================
	s := &http.Server{
		Addr:         ":" + port,
		Handler:      e,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Fatal(e.StartServer(s))
}
