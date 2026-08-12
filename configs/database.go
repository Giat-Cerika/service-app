package configs

import (
	"context"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		PrepareStmt: false,
		Logger:      logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("Gagal koneksi ke database: %v", err)
	}

	// ================================================================
	// CONNECTION POOL — disesuaikan untuk Biznet NEO Lite MS4.2
	//   2 vCPU  •  4 GB RAM  •  60 GB SSD
	//
	// Aturan praktis:
	//   MaxOpenConns  = 2..3× jumlah CPU → 30  (cukup untuk 300 user
	//                   karena tiap HTTP req hanya pakai 1 koneksi
	//                   selama < 50ms, selebihnya antre di pool)
	//   MaxIdleConns  = ~⅓ MaxOpen → 10  (hemat RAM, 1 koneksi idle
	//                   ≈ 2 MB di Postgres)
	//   ConnMaxLife   = 10 menit   (hindari koneksi mati / stale)
	//   ConnMaxIdle   = 3 menit    (lepas koneksi idle lebih cepat)
	// ================================================================
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(30)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(10 * time.Minute)
		sqlDB.SetConnMaxIdleTime(3 * time.Minute)
	}

	log.Println("✅ Database connected — pool: MaxOpen=30 MaxIdle=10")
	DB = db
	return db
}

// GetDB mengembalikan *gorm.DB global untuk keperluan transaction di service layer.
func GetDB() *gorm.DB {
	return DB
}

// RunTransaction adalah helper untuk menjalankan fungsi di dalam DB transaction.
// Jika fn mengembalikan error, transaction di-rollback secara otomatis.
// Jika fn berhasil (return nil), transaction di-commit.
func RunTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return DB.WithContext(ctx).Transaction(fn)
}
