package middleware

import (
	"net/http"
	"testing"

	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/util/testutil"
)

func TestLogging(t *testing.T) {
	tests := []testutil.LoggingTestDefinition{
		{
			Name: "Logging 200 OK response",
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			ExpectedStatusCode: http.StatusOK,
			ExpectedMethod:     http.MethodGet,
			ExpectedPath:       "/test-path",
		},
		{
			Name: "Logging 404 Not Found response",
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			ExpectedStatusCode: http.StatusNotFound,
			ExpectedMethod:     http.MethodGet,
			ExpectedPath:       "/not-found-path",
		},
		{
			Name: "Logging 500 Internal Server Error response",
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			ExpectedStatusCode: http.StatusInternalServerError,
			ExpectedMethod:     http.MethodPost,
			ExpectedPath:       "/error-path",
		},
	}

	testutil.RunLoggerTests(t, Logging, tests)
}