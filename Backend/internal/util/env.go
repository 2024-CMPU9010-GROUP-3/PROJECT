package util

import (
	"errors"
	"log"
	"os"
)

const (
	envPort           = "MAGPIE_PORT"
	envDbUrl          = "MAGPIE_DB_URL"
	envMigrationsPath = "MAGPIE_DB_MIGRATIONS"
	envCorsOrigins    = "MAGPIE_CORS_ALLOWED_ORIGINS"
	envCorsMethods    = "MAGPIE_CORS_ALLOWED_METHODS"
	envJwtSecret      = "MAGPIE_JWT_SECRET"
	envJwtExpiry      = "MAGPIE_JWT_EXPIRY"
	envLogin          = "LOGIN"
	envPassword       = "PASSWORD"
	envHost           = "HOST"
	envDbName         = "DATABASE_NAME"
)

var (
	defaults = map[string]string{
		envPort:           "8080",
		envDbUrl:          "",
		envMigrationsPath: "./sql/migratios",
		envCorsOrigins:    "",
		envCorsMethods:    "",
		envJwtSecret:      "",
		envJwtExpiry:      "",
		envLogin:          "",
		envPassword:       "",
		envHost:           "",
		envDbName:         "",
	}
)

func Get(name string) (string, bool) {
	env, set := os.LookupEnv(name)
	if !set && len(defaults[name]) != 0 {
		log.Printf("Warning: %s not found in environment, falling back to default", name)
		return defaults[name], true
	}
	return env, set
}

func Validate() error {
	_, dbUrlSet := Get(envDbUrl)
	_, loginSet := Get(envLogin)
	_, passwordSet := Get(envPassword)
	_, hostSet := Get(envHost)
	_, nameSet := Get(envDbName)

	if !dbUrlSet && (!loginSet || !passwordSet || !hostSet || !nameSet) {
		return errors.New("Error: No database URL or database login credentials found in environment")
	}

	_, jwtSecretSet := Get(envJwtSecret)
	if !jwtSecretSet {
		return errors.New("Error: No JWT Secret found in environment")
	}

	return nil

}
