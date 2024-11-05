package util

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type TokenKey string

type placeholder struct {
	IsPlaceholder bool
	Endpoint      string
}

func Placeholder(endpoint string) *placeholder {
	return &placeholder{true, endpoint}
}

func CheckResponseError(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, fmt.Sprintf("response error, %v", err), http.StatusInternalServerError)
	}
}

func GetLatestMigrationVersion(migrationsPath string) (uint, error) {
	var latestVersion uint // use uint to match migrate package

	err := filepath.Walk(migrationsPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
					return err
			}
			// only look at up-migration files
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".up.sql") {
					parts := strings.Split(info.Name(), "_")
					if len(parts) > 0 {
							version, err := strconv.Atoi(parts[0]) // migration version is the first component
							if err == nil && uint(version) > latestVersion {
									latestVersion = uint(version)
							}
					}
			}
			return nil // no error occurred
	})

	return latestVersion, err
}