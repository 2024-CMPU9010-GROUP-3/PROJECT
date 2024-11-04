package testutil

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/middleware"
)

type MiddlewareTestDefinition struct {
	Name               string
	IdPathParam        string
	TokenUserID        string
	ExpectedStatusCode int
	ExpectedBody       string
}

func executeMiddlewareTest(t *testing.T, test MiddlewareTestDefinition, middleware middleware.Middleware) {
	req := httptest.NewRequest(http.MethodGet, "/some-path", nil)

	req = req.WithContext(context.WithValue(req.Context(), "token_user_id", test.TokenUserID))
	req = setPathValue(req, "id", test.IdPathParam)

	rr := httptest.NewRecorder()

	handler := middleware(http.HandlerFunc(mockNextHandler))

	handler.ServeHTTP(rr, req)

	if rr.Code != test.ExpectedStatusCode {
		t.Errorf("expected status %d, got %d", test.ExpectedStatusCode, rr.Code)
	}

	if !strings.Contains(rr.Body.String(), test.ExpectedBody) {
		t.Errorf("expected response body to contain %q, got %q", test.ExpectedBody, rr.Body.String())
	}
}

func RunMiddlewareTest(t *testing.T, middleware middleware.Middleware, tests []MiddlewareTestDefinition) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			executeMiddlewareTest(t, tt, middleware)
		})
	}
}

func setPathValue(req *http.Request, key, value string) *http.Request {
	ctx := context.WithValue(req.Context(), key, value)
	return req.WithContext(ctx)
}

func mockNextHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("next handler called"))
}
