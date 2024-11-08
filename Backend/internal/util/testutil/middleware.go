package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/util"
)

type MiddlewareTestDefinition struct {
	Name                 string
	IdPathParam          string
	TokenUserId          interface{}
	ExpectedStatusCode   int
	ExpectedBody         string
	ExpectedBodyContains string
	AuthCookieValue      string
	EnvJwtSecret         string
}

type LoggingTestDefinition struct {
	Name               string
	HandlerFunc        http.HandlerFunc
	ExpectedStatusCode int
	ExpectedMethod     string
	ExpectedPath       string
}

func executeMiddlewareTest(t *testing.T, test MiddlewareTestDefinition, middleware func(http.Handler) http.Handler) {
	req := httptest.NewRequest(http.MethodGet, "/some-path", nil)

	req = req.WithContext(context.WithValue(req.Context(), util.TokenKey("token_user_id"), test.TokenUserId))
	req.SetPathValue("id", test.IdPathParam)

	if test.AuthCookieValue != "" {
		req.AddCookie(&http.Cookie{
			Name:  "magpie_auth",
			Value: test.AuthCookieValue,
		})
	}

	if test.EnvJwtSecret != "" {
		t.Setenv("MAGPIE_JWT_SECRET", test.EnvJwtSecret)
	}

	rr := httptest.NewRecorder()

	handler := middleware(http.HandlerFunc(mockNextHandler))

	handler.ServeHTTP(rr, req)

	if rr.Code != test.ExpectedStatusCode {
		t.Errorf("expected status %d, got %d", test.ExpectedStatusCode, rr.Code)
	}

	if test.ExpectedBody != "" {
		compactedJson := &bytes.Buffer{}
		err := json.Compact(compactedJson, []byte(test.ExpectedBody))
		if err != nil {
			t.Errorf("could not flatten expected JSON, this is due to incorrect test case definition")
		}

		// this is needed because the response body always includes a newline
		compactedJson.WriteByte(0x0a)

		if rr.Body.String() != compactedJson.String() {
			t.Errorf("expected JSON output %s, got %s", compactedJson.String(), rr.Body.String())
		}
	}
	if test.ExpectedBodyContains != "" && !strings.Contains(rr.Body.String(), test.ExpectedBodyContains) {
		t.Errorf("expected body to contain %s, got %s", test.ExpectedBodyContains, rr.Body.String())
	}
}

func executeLoggingTest(t *testing.T, test LoggingTestDefinition, loggingMiddleware func(http.Handler) http.Handler) {
	t.Run(test.Name, func(t *testing.T) {
		logOutput := captureLogs(func() {
			req := httptest.NewRequest(test.ExpectedMethod, test.ExpectedPath, nil)
			rr := httptest.NewRecorder()

			handler := loggingMiddleware(http.HandlerFunc(test.HandlerFunc))
			handler.ServeHTTP(rr, req)
		})

		if !strings.Contains(logOutput, strconv.Itoa(test.ExpectedStatusCode)) {
			t.Errorf("expected log to contain status code %d (%s), but got: %s",
				test.ExpectedStatusCode, http.StatusText(test.ExpectedStatusCode), logOutput)
		}
		if !strings.Contains(logOutput, test.ExpectedMethod) {
			t.Errorf("expected log to contain method %s, but got: %s", test.ExpectedMethod, logOutput)
		}
		if !strings.Contains(logOutput, test.ExpectedPath) {
			t.Errorf("expected log to contain path %s, but got: %s", test.ExpectedPath, logOutput)
		}

		if !strings.Contains(logOutput, "ms") && !strings.Contains(logOutput, "s") {
			t.Error("expected log to contain duration, but no duration was found")
		}
	})
}

func RunLoggerTests(t *testing.T, loggingMiddleware func(http.Handler) http.Handler, tests []LoggingTestDefinition) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			executeLoggingTest(t, tt, loggingMiddleware)
		})
	}
}

func RunMiddlewareTests(t *testing.T, middleware func(http.Handler) http.Handler, tests []MiddlewareTestDefinition) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			executeMiddlewareTest(t, tt, middleware)
		})
	}
}

func mockNextHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("next handler called"))
}

func captureLogs(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil) // Restore original log output
	f()
	return buf.String()
}
