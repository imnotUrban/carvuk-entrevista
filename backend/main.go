package main

import (
	"fmt"
	"log"
	"os"

	"backend/delivery"
	boilerplateentity "backend/entity/boilerplate"
	boletaentity "backend/entity/boleta"
	compraentity "backend/entity/compra"
	productoentity "backend/entity/producto"
	usuarioentity "backend/entity/usuario"
	boletarepo "backend/repository/boleta"
	productorepo "backend/repository/producto"
	boletaservice "backend/service/boleta"
	productoservice "backend/service/producto"

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

	if err := db.AutoMigrate(&boilerplateentity.Boilerplate{}, &productoentity.Producto{},
		&usuarioentity.Usuario{}, &compraentity.Compra{}, &compraentity.CompraItem{}, &boletaentity.Boleta{}); err != nil {
		log.Fatalf("failed to run auto migrations: %v", err)
	}

	if err := productoservice.NewService(productorepo.NewRepository(db)).SeedIfEmpty(); err != nil {
		log.Fatalf("failed to seed productos: %v", err)
	}

	if err := boletaservice.NewService(boletarepo.NewRepository(db), productorepo.NewRepository(db)).SeedUsuarioIfEmpty(); err != nil {
		log.Fatalf("failed to seed usuario: %v", err)
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
