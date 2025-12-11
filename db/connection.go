package db

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(envPath string) *gorm.DB {
	err := godotenv.Load(envPath)

	if err != nil {
		log.Fatal("Error loading environment file", err)
	}

	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))

	if err != nil {
		fmt.Println("Invalid POSTGRES_PORT:", os.Getenv("POSTGRES_PORT"))
	}

	var dsn string

	env := os.Getenv("APP_ENV")

	switch env {
	case "prod", "dev":
		dsn = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("POSTGRES_HOST"),
			port,
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_DB"),
		)
	case "test":
		dsn = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("POSTGRES_HOST"),
			port,
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_TEST_DB"),
		)
	default:
		log.Fatalf("unsupported environment: %s", env)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		panic("failed to connect database")
	}

	return db
}
