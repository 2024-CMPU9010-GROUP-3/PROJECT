package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/handlers"
	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/middleware"
	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/routes"
	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/util"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	pgxgeom "github.com/twpayne/pgx-geom"
)

const defaultPort = "8080"
const portEnv = "MAGPIE_PORT"
const dbUrlEnv = "MAGPIE_DB_URL"
const corsAllowedOriginsEnv = "MAGPIE_CORS_ALLOWED_ORIGINS"
const corsAllowedMethodsEnv = "MAGPIE_CORS_ALLOWED_METHODS"

const migrationsPath = "./sql/migrations"

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

	allowedOriginsEnv := os.Getenv(corsAllowedOriginsEnv)
	if len(allowedOriginsEnv) == 0 {
		log.Printf("Warning: Could not find allowed origins for CORS in environment variables")
	}
	allowedMethodsEnv := os.Getenv(corsAllowedMethodsEnv)
	if len(allowedMethodsEnv) == 0 {
		log.Printf("Warning: Could not find allowed methods for CORS in environment variables")
	}

	allowedOrigins := strings.Split(allowedOriginsEnv, " ")
	allowedMethods := strings.Split(allowedMethodsEnv, " ")

	corsManager := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   allowedMethods,
		AllowCredentials: true, // needed for login
	})

	middlewares := middleware.CreateMiddlewareStack(
		middleware.Logging,
		corsManager.Handler,
	)
	router := routes.Router

	ctx := context.Background()

	m, err := migrate.New(fmt.Sprintf("file://%s", migrationsPath), dbUrl)
	if err != nil {
		log.Fatalf("Could not create database migration client: %v\n", err)
		os.Exit(1)
	}

	currentMigration, dirty, err := m.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			log.Printf("No migration version found in database\n")
		} else {
			log.Fatalf("Could not read migration version from database: %v\n", err)
			os.Exit(1)
		}
	}

	if dirty {
		log.Printf("WARNING: The database migration state is currently considered dirty")
	}

	log.Printf("Checking for new database migrations")
	latestMigration, err := util.GetLatestMigrationVersion(migrationsPath)

	if latestMigration > currentMigration {
		log.Printf("Found newer database migrations (current: V%d, latest: V%d), attemping upgrade...", currentMigration, latestMigration)
		err = m.Up()
		if err != nil {
			log.Fatalf("An error occurred during database migration: %v", err)
			os.Exit(1)
		}
		log.Printf("Successfully migrated database: V%d => V%d", currentMigration, latestMigration)
	} else {
		log.Printf("Database is up to date")
	}

	
	dbConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Fatalf("Could not parse pool config: %+v\n", err)
		os.Exit(1)
	}
	
	initialDbPool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		log.Fatalf("Could not create database pool: %+v\n", err)
		os.Exit(1)
	}
	customTypes, err := util.GetCustomDataTypes(ctx, initialDbPool)
	if err != nil {
		log.Fatalf("Could not get custom types from database pool: %+v\n", err)
		os.Exit(1)
	}

	dbConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		for _, t := range customTypes {
			conn.TypeMap().RegisterType(t)
		}

		err = pgxgeom.Register(ctx, conn)
		if err != nil {
			log.Fatalf("Could not register custom datatypes: %v\n", err)
		}
		return err
	}
	initialDbPool.Close() // close and reopen

	dbPool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		log.Fatalf("Could not connect to database: %v\n", err)
		os.Exit(1)
	} else {
		log.Println("Successfully connected to database")
	}
	defer dbPool.Close()

	handlers.RegisterDatabaseConnection(&ctx, dbPool)

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
