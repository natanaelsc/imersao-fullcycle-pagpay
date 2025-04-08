package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/natanaelsc/imersao-fullcycle-pagpay/gateway-api/internal/repository"
	"github.com/natanaelsc/imersao-fullcycle-pagpay/gateway-api/internal/service"
	"github.com/natanaelsc/imersao-fullcycle-pagpay/gateway-api/internal/web/server"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Carrega variáveis de ambiente do arquivo .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Configura conexão com PostgreSQL usando variáveis de ambiente
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "db"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "gateway"),
		getEnv("DB_SSL_MODE", "disable"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	accountRepository := repository.NewAccountRepository(db)
	accountService := service.NewAccountService(accountRepository)

	port := getEnv("PORT", "8080")
	srv := server.NewServer(accountService, port)
	srv.ConfigureRoutes()

	srv.Start()
	if err := srv.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
