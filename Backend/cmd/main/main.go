package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/handlers"
	"github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/middleware"
	"github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/routes"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/twpayne/go-geos"
	pgxgeos "github.com/twpayne/pgx-geos"
)

const defaultPort = "8080"
const portEnv = "MAGPIE_PORT"
const dbUrlEnv = "MAGPIE_DB_URL"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found")
	}

	port := os.Getenv(portEnv)
	if len(port) == 0 {
		port = defaultPort
	}

	dbUrl := os.Getenv(dbUrlEnv)
	if len(dbUrl) == 0 {
		// Try to build the connection string from individual components in case we are running in Kubernetes
		log.Printf("Could not find connection string in environment variables, trying to build from components")
		dbUser := os.Getenv("LOGIN")
		dbPassword := os.Getenv("PASSWORD")
		dbHost := os.Getenv("HOST")
		dbName := os.Getenv("DATABASE_NAME")

		if len(dbUser) == 0 || len(dbPassword) == 0 || len(dbHost) == 0 || len(dbName) == 0 {
			log.Fatal("Database connection parameters are missing in environment variables")
			os.Exit(1)
		}

		dbUrl = fmt.Sprintf("postgresql://%s:%s@%s/%s", dbUser, dbPassword, dbHost, dbName)
		log.Printf("%s\n", dbUrl)
	}

	router := routes.Router

	middlewares := middleware.CreateMiddlewareStack(
		middleware.Logging,
	)

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		log.Fatalf("Could not connect to database: %v\n", err)
		os.Exit(1)
	} else {
		log.Println("Successfully connected to database")
	}
	defer conn.Close(ctx)

	err = pgxgeos.Register(ctx, conn, geos.NewContext())
	if err != nil {
		log.Fatalf("Could not register geo datatype: %v\n", err)
		os.Exit(1)
	}

	handlers.RegisterDatabaseConnection(&ctx, conn)

	server := http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: middlewares(router),
	}

	log.Printf("Listening on port %v\n", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Printf("Server: %v\n", err)
	}
}
