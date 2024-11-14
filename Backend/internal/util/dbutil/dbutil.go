package dbutil

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/util"
	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/util/env"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxgeom "github.com/twpayne/pgx-geom"
)

func SetupDatabase() (*pgxpool.Pool, error) {
	dbUrl := env.GetDbUrl()
	dbRetries, _ := env.GetInt(env.EnvDbRetries)
	dbRetryInterval, _ := env.GetInt(env.EnvDbRetryInterval)
	migrationsPath, _ := env.Get(env.EnvMigrationsPath)
	
	dbConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		return nil, fmt.Errorf("Could not parse database config: %s", err)
	}
	
	err = awaitDatabase(dbConfig, dbRetries, dbRetryInterval)
	if err != nil {
		return nil, fmt.Errorf("Could not establish database connection: %s", err)
	}

	err = runMigrations(dbConfig, migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("Could not run database migrations: %s", err)
	}

	dbConfig, err = registerCustomDatatypes(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("Could not register custom datatypes: %s", err)
	}

	return pgxpool.NewWithConfig(context.Background(), dbConfig)
}

func awaitDatabase(dbConfig *pgxpool.Config, dbRetries int, dbRetryInterval int) error {

	dbConfig.ConnConfig.ConnectTimeout = time.Duration(dbRetryInterval  * int(time.Second))

	log.Printf("Attempting to connect to database")
	_, err := pgx.ConnectConfig(context.Background(), dbConfig.ConnConfig)
	if err != nil {
		log.Printf("Waiting for database...")
	}
	for err != nil && dbRetries > 0 {
		_, err = pgx.ConnectConfig(context.Background(), dbConfig.ConnConfig)
		dbRetries--
		log.Printf("Waiting for database (Retries left: %d)", dbRetries)
	}
	return err
}

func runMigrations(dbConfig *pgxpool.Config, migrationsPath string) error {
	migrateClient, err := migrate.New(fmt.Sprintf("file://%s", migrationsPath), dbConfig.ConnString())
	if err != nil {
		return fmt.Errorf("Could not create database migration client: %v\n", err)
	}

	currentMigration, dirty, err := migrateClient.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			log.Printf("No migration version found in database\n")
		} else {
			return fmt.Errorf("Could not read migration version from database: %v\n", err)
		}
	}

	if dirty {
		log.Printf("WARNING: The database migration state is currently considered dirty")
	}

	log.Printf("Checking for new database migrations")
	latestMigration, err := util.GetLatestMigrationVersion(migrationsPath)
	if err != nil {
		return fmt.Errorf("An error occurred while searching for newest migration version: %v", err)
	}

	if latestMigration > currentMigration {
		log.Printf("Found newer database migrations (current: V%d, latest: V%d), attemping upgrade...", currentMigration, latestMigration)
		err = migrateClient.Up()
		if err != nil {
			return fmt.Errorf("An error occurred during database migration: %v", err)
		}
		log.Printf("Successfully migrated database: V%d => V%d", currentMigration, latestMigration)
	} else {
		log.Printf("Database is up to date")
	}
	return nil
}

func registerCustomDatatypes(dbConfig *pgxpool.Config) (*pgxpool.Config, error) {
	dbPool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		return nil, fmt.Errorf("Could not connect to database to register custom types: %v\n", err)
	}

	customTypes, err := util.GetCustomDataTypes(context.Background(), dbPool)
	if err != nil {
		return nil, fmt.Errorf("Could not get custom types from database pool: %+v\n", err)
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
	return dbConfig.Copy(), nil
}
