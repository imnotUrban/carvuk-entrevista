package main

import (
	"fmt"
	"log"
	"os"

	"backend/delivery"
	boilerplateentity "backend/entity/boilerplate"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on environment variables")
	}

	db, err := connectDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&boilerplateentity.Boilerplate{}); err != nil {
		log.Fatalf("failed to run auto migrations: %v", err)
	}

	// AutoMigrate cannot express a partial unique index, so it's applied
	// separately here. It mirrors migrations/0001_create_boilerplate.up.sql and
	// keeps a soft-deleted row's code free for reuse.
	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_boilerplate_code_active
		ON boilerplate (code)
		WHERE deleted_at IS NULL
	`).Error; err != nil {
		log.Fatalf("failed to create partial unique index on boilerplate.code: %v", err)
	}

	router := delivery.NewRouter(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server listening on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func connectDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "carvuk"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_SSLMODE", "disable"),
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
