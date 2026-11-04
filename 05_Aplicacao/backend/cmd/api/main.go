package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"legisanalytica/internal/activity"
	"legisanalytica/internal/expense"
	"legisanalytica/internal/graphql"
	"legisanalytica/internal/politician"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func connectPostgres() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\n", log.LstdFlags), logger.Config{
			SlowThreshold: 2 * time.Minute,
			LogLevel:      logger.Warn,
		}),
	})
}

func initializeApiHandler(gormDB *gorm.DB) *graphql.GraphQLHandler {
	pRepo := politician.NewRepository(gormDB)
	pSvc := politician.NewService(pRepo)

	dRepo := expense.NewRepository(gormDB)
	dSvc := expense.NewService(dRepo)

	aSvcRepo := activity.NewRepository(gormDB)
	aSvc := activity.NewService(aSvcRepo)

	return graphql.NewGraphQLHandler(pSvc, dSvc, aSvc)
}

const SERVER_PORT = ":8080"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Aviso: Erro ao carregar arquivo .env: %v", err)
	}

	gormDB, err := connectPostgres()
	if err != nil {
		log.Fatalf("Erro ao conectar ao Postgres: %v", err)
	}

	handler := initializeApiHandler(gormDB)

	mux := http.NewServeMux()
	mux.Handle("/graphql", withCORS(os.Getenv("CORS_ALLOWED_ORIGIN"), http.HandlerFunc(handler.HandleGraphQL)))

	srv := &http.Server{
		Addr:              SERVER_PORT,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      70 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 16,
	}

	fmt.Printf("Server running on %s (HTTP)\n", SERVER_PORT)
	log.Fatal(srv.ListenAndServe())
}
