package env

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	EnvPort            = "MAGPIE_PORT"
	EnvDbUrl           = "MAGPIE_DB_URL"
	EnvMigrationsPath  = "MAGPIE_DB_MIGRATIONS"
	EnvCorsOrigins     = "MAGPIE_CORS_ALLOWED_ORIGINS"
	EnvCorsMethods     = "MAGPIE_CORS_ALLOWED_METHODS"
	EnvJwtSecret       = "MAGPIE_JWT_SECRET"
	EnvJwtExpiry       = "MAGPIE_JWT_EXPIRY"
	EnvLogin           = "LOGIN"
	EnvPassword        = "PASSWORD"
	EnvHost            = "HOST"
	EnvDbName          = "DATABASE_NAME"
	EnvDbRetries       = "MAGPIE_DB_RETRIES"
	EnvDbRetryInterval = "MAGPIE_DB_RETRY_INTERVAL"
	EnvShutdownTimeout = "MAGPIE_SHUTDOWN_TIMEOUT"
)

var (
	defaults = map[string]string{
		EnvPort:            "8080",
		EnvDbUrl:           "",
		EnvMigrationsPath:  "./sql/migrations",
		EnvCorsOrigins:     "",
		EnvCorsMethods:     "",
		EnvJwtSecret:       "",
		EnvJwtExpiry:       "168h",
		EnvLogin:           "",
		EnvPassword:        "",
		EnvHost:            "",
		EnvDbName:          "",
		EnvDbRetries:       "6",
		EnvDbRetryInterval: "10",
		EnvShutdownTimeout: "10",
	}
)

func Get(name string) (string, bool) {
	env, set := os.LookupEnv(name)
	if !set {
		if len(defaults[name]) != 0 {
			return defaults[name], true
		} else {
			log.Printf("Warning: %s not found in environment", name)
			return "", false
		}
	}
	return env, set
}

func GetInt(name string) (int, bool) {
	env, _ := Get(name)
	envInt, err := strconv.Atoi(env)
	if err != nil {
		return 0, false
	}
	return envInt, true
}

func Load() {
	loadErr := godotenv.Load()
	if loadErr != nil {
		fmt.Printf("Warning: No .env file found")
	}
}

func Validate() error {
	_, dbUrlSet := Get(EnvDbUrl)
	_, loginSet := Get(EnvLogin)
	_, passwordSet := Get(EnvPassword)
	_, hostSet := Get(EnvHost)
	_, nameSet := Get(EnvDbName)

	if !dbUrlSet && (!loginSet || !passwordSet || !hostSet || !nameSet) {
		return errors.New("Error: No database URL or database login credentials found in environment")
	}

	_, jwtSecretSet := Get(EnvJwtSecret)
	if !jwtSecretSet {
		return errors.New("Error: No JWT Secret found in environment")
	}

	_, jwtExpirySet := Get(EnvJwtExpiry)
	if !jwtExpirySet {
		return errors.New("Error: No JWT Expiry found in environment")
	}

	_, migrationsPathSet := Get(EnvMigrationsPath)
	if !migrationsPathSet {
		return errors.New("Error: No migrations path found in environment")
	}

	_, dbRetriesSet := GetInt(EnvDbRetries)
	if !dbRetriesSet {
		return errors.New("Error: No retry amount for database connection found in environment")
	}

	_, dbRetryIntervalSet := GetInt(EnvDbRetryInterval)
	if !dbRetryIntervalSet {
		return errors.New("Error: No retry interval for database connection found in environment")
	}

	return nil

}

func GetDbUrl() string {
	dbUrl, dbUrlSet := Get(EnvDbUrl)
	if dbUrlSet {
		return dbUrl
	} else {
		log.Printf("Could not find connection string in environment variables, trying to build from components")
		dbUser := os.Getenv("LOGIN")
		dbPassword := os.Getenv("PASSWORD")
		dbHost := os.Getenv("HOST")
		dbName := os.Getenv("DATABASE_NAME")

		return fmt.Sprintf("postgresql://%s:%s@%s/%s", dbUser, dbPassword, dbHost, dbName)
	}
}
