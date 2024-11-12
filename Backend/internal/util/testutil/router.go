package testutil

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
)

type RouterTestDefinition struct {
	Name               string
	Path               string
	Method             string
	ExpectedStatusCode int
	MockSetup          func(mock pgxmock.PgxPoolIface)
}

func executeRouterTest(t *testing.T, tt RouterTestDefinition, router *http.ServeMux, mock pgxmock.PgxPoolIface) {
	tt.MockSetup(mock)
	req := httptest.NewRequest(tt.Method, tt.Path, nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != tt.ExpectedStatusCode {
		t.Errorf("Expected status code %d, got %d - for route %s", tt.ExpectedStatusCode, rr.Code, tt.Path)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func RunRouterTests(t *testing.T, tests []RouterTestDefinition, router *http.ServeMux, mock pgxmock.PgxPoolIface) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			executeRouterTest(t, tt, router, mock)
		})
	}
}
