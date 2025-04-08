package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/devfullcycle/imersao22/go-gateway/internal/repository"
	"github.com/devfullcycle/imersao22/go-gateway/internal/service"
	"github.com/devfullcycle/imersao22/go-gateway/internal/web/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Getenv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbConnectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		Getenv("DB_HOST", "localhost"),
		Getenv("DB_PORT", "5432"),
		Getenv("DB_USER", "postgres"),
		Getenv("DB_PASSWORD", "postgres"),
		Getenv("DB_NAME", "gateway"),
		Getenv("DB_SSL_MODE", "disable"),
	)

	db, err := sql.Open("postgres", dbConnectionString)

	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging database: ", err)
	}

	accountRepository := repository.NewAccountRepository(db)
	accountService := service.NewAccountService(accountRepository)

	server := server.NewServer(accountService, Getenv("PORT", "8080"))
	server.ConfigureRoutes()

	log.Println("Server started on port", Getenv("PORT", "8080"))

	if err := server.Start(); err != nil {
		log.Fatal("Error starting server: ", err)
	}

}
