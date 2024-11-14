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
	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/util/env"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	pgxgeom "github.com/twpayne/pgx-geom"
)

func main() {
	env.Load()

	err := env.Validate()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	port, _ := env.Get(env.EnvPort)

	dbUrl := env.GetDbUrl()

	allowedOriginsEnv, _ := env.Get(env.EnvCorsOrigins)
	allowedMethodsEnv, _ := env.Get(env.EnvCorsMethods)

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

	migrationsPath, _ := env.Get(env.EnvMigrationsPath)

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
	if err != nil {
		log.Fatalf("An error occurred while searching for newest migration version: %v", err)
		os.Exit(1)
	}

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
