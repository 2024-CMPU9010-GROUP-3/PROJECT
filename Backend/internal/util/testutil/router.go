package testutil

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type RouterTestDefinition struct {
	Name               string
	Path               string
	Method             string
	ExpectedStatusCode int
}

func executeRouterTest(t *testing.T, tt RouterTestDefinition, router *http.ServeMux) {
	req := httptest.NewRequest(tt.Method, tt.Path, nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != tt.ExpectedStatusCode {
		t.Errorf("Expected status code %d, got %d - for route %s", tt.ExpectedStatusCode, rr.Code, tt.Path)
	}
}

func RunRouterTests(t *testing.T, tests []RouterTestDefinition, router *http.ServeMux) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			executeRouterTest(t, tt, router)
		})
	}
}
