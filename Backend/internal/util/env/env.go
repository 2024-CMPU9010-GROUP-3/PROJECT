package env

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	EnvPort           = "MAGPIE_PORT"
	EnvDbUrl          = "MAGPIE_DB_URL"
	EnvMigrationsPath = "MAGPIE_DB_MIGRATIONS"
	EnvCorsOrigins    = "MAGPIE_CORS_ALLOWED_ORIGINS"
	EnvCorsMethods    = "MAGPIE_CORS_ALLOWED_METHODS"
	EnvJwtSecret      = "MAGPIE_JWT_SECRET"
	EnvJwtExpiry      = "MAGPIE_JWT_EXPIRY"
	EnvLogin          = "LOGIN"
	EnvPassword       = "PASSWORD"
	EnvHost           = "HOST"
	EnvDbName         = "DATABASE_NAME"
)

var (
	defaults = map[string]string{
		EnvPort:           "8080",
		EnvDbUrl:          "",
		EnvMigrationsPath: "./sql/migrations",
		EnvCorsOrigins:    "",
		EnvCorsMethods:    "",
		EnvJwtSecret:      "",
		EnvJwtExpiry:      "",
		EnvLogin:          "",
		EnvPassword:       "",
		EnvHost:           "",
		EnvDbName:         "",
	}
)

func Get(name string) (string, bool) {
	env, set := os.LookupEnv(name)
	if !set {
		if len(defaults[name]) != 0 {
			log.Printf("Warning: %s not found in environment, falling back to default", name)
			return defaults[name], true
		} else {
			log.Printf("Warning: %s not found in environment", name)
			return "", false
		}
	}
	return env, set
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