package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/dtos"
	"github.com/pashagolub/pgxmock/v4"
)

type HandlerTestDefinition struct {
	Name            string
	Method          string
	Route           string
	InputJSON       string
	MockSetup       func(mock pgxmock.PgxPoolIface)
	ExpectedStatus  int
	ExpectedError   string
	ExpectedJSON    string
	PathParams      map[string]string
	QueryParams     map[string]string
	Env             map[string]string
}

var jwtPattern = regexp.MustCompile(`^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+$`)

func executeHandlerTest(t *testing.T, tt HandlerTestDefinition, handlerFunc func(rr http.ResponseWriter, req *http.Request), mock pgxmock.PgxPoolIface) {
	tt.MockSetup(mock)

	for k, v := range tt.Env {
		t.Setenv(k, v)
	}

	req, err := http.NewRequest(tt.Method, tt.Route, bytes.NewBuffer([]byte(tt.InputJSON)))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range tt.PathParams {
		req.SetPathValue(k, v)
	}

	q := req.URL.Query()
	for k, v := range tt.QueryParams {
		q.Add(k, v)
	}

	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()

	// FUNCTION OF INTEREST
	handlerFunc(rr, req)

	if status := rr.Code; status != tt.ExpectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v", status, tt.ExpectedStatus)
	}

	if tt.ExpectedError != "" {
		var responseBody dtos.ResponseDto
		if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if responseBody.Error.ErrorMsg != tt.ExpectedError {
			t.Errorf("expected error message \"%v\", got \"%v\"", tt.ExpectedError, responseBody.Error.ErrorMsg)
		}
	}

	if tt.ExpectedJSON != "" {
		compactedJson := &bytes.Buffer{}
		err := json.Compact(compactedJson, []byte(tt.ExpectedJSON))
		if err != nil {
			t.Errorf("could not flatten expected JSON, this is due to incorrect test case definition")
		}

		// this is needed because the response body always includes a newline
		compactedJson.WriteByte(0x0a)

		if rr.Body.String() != compactedJson.String() {
			t.Errorf("expected JSON output %s, got %s", compactedJson.String(), rr.Body.String())
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func RunHandlerTests(t *testing.T, handlerFunc func(rr http.ResponseWriter, req *http.Request), mock pgxmock.PgxPoolIface, tests []HandlerTestDefinition) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			executeHandlerTest(t, tt, handlerFunc, mock)
		})
	}
}
