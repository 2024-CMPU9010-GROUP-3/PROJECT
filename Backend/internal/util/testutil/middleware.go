package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MiddlewareTestDefinition struct {
	Name                 string
	IdPathParam          string
	TokenUserId          interface{}
	ExpectedStatusCode   int
	ExpectedBody         string
	ExpectedBodyContains string
}

type tokenKey string

func executeMiddlewareTest(t *testing.T, test MiddlewareTestDefinition, middleware func(http.Handler) http.Handler) {
	req := httptest.NewRequest(http.MethodGet, "/some-path", nil)

	req = req.WithContext(context.WithValue(req.Context(), "token_user_id", test.TokenUserId))
	req.SetPathValue("id", test.IdPathParam)

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

func RunMiddlewareTests(t *testing.T, middleware func(http.Handler) http.Handler, tests []MiddlewareTestDefinition) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			executeMiddlewareTest(t, tt, middleware)
		})
	}
}

func mockNextHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("next handler called"))
}
