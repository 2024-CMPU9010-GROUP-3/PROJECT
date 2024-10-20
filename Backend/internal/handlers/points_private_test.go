//go:build private

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	db "github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/db/private"
	customErrors "github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/errors"
	resp "github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/responses"
	"github.com/pashagolub/pgxmock/v4"
)

// Struct for test cases
type PointsHandlerTest struct {
	name           string
	inputJSON      string
	mockSetup      func(mock pgxmock.PgxPoolIface)
	expectedStatus int
	expectedError  string
	expectedJSON   string
}

func runHandlerTest(t *testing.T, tt PointsHandlerTest, handlerFunc func(rr http.ResponseWriter, req *http.Request), mock pgxmock.PgxPoolIface) {
	tt.mockSetup(mock)

	req, err := http.NewRequest("POST", "/points", bytes.NewBuffer([]byte(tt.inputJSON)))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// FUNCTION OF INTEREST
	handlerFunc(rr, req)

	if status := rr.Code; status != tt.expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
	}

	if tt.expectedError != "" {
		var responseBody resp.ResponseDto
		if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		t.Logf("%s", rr.Body.Bytes())
		t.Logf("%+v", responseBody)
		if responseBody.Error.ErrorMsg != tt.expectedError {
			t.Errorf("expected error message %v, got %v", tt.expectedError, responseBody.Error.ErrorMsg)
		}
	}

	if tt.expectedJSON != "" {
		if rr.Body.String() != tt.expectedJSON {
			t.Errorf("expected JSON output %v, got %v", tt.expectedJSON, rr.Body.String())
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestPointsHandlerHandlePost(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := context.Background()

	RegisterDatabaseConnection(&ctx, mock)

	handler := &PointsHandler{}

	tests := []PointsHandlerTest{
		{
			name: "Valid input",
			inputJSON: `{
				"longlat": {
					"type": "Point",
					"coordinates": [11, 12]
				},
				"type": "placeholder1",
				"details": {
					"test": 1234
				}
			}`,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("INSERT INTO points").
					WithArgs(
						pgxmock.AnyArg(),
						db.PointType("placeholder1"),
						[]byte(`{"test":1234}`),
					).WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(int64(1)))
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid input",
			inputJSON: `{
				"longlat": {
					"type": "Point",
					"coordinates": [11, 12]
				},
				"type1": "placeholder1",
				"details": {
					"test": 1234
				}
			}`,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				// No mock setup needed, handler should return error before making db request
			},
			expectedStatus: http.StatusBadRequest,
			expectedError: customErrors.Payload.InvalidPayloadPointError.ErrorMsg,
		},
		{
			name: "Invalid geometry",
			inputJSON: `{
				"longlat": {
					"type": "InvalidType",
					"coordinates": [11, 12]
				},
				"type": "placeholder1",
				"details": {
					"test": 1234
				}
			}`,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				// No mock setup needed, handler should return error before making db request
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  customErrors.Payload.InvalidPayloadPointError.ErrorMsg,
		},
		{
			name: "Valid geometry, but not point",
			inputJSON: `{
				"longlat": {
					"type": "Polygon",
					"coordinates": [
						[
							[100.0, 0.0],
							[101.0, 0.0],
							[101.0, 1.0],
							[100.0, 1.0],
							[100.0, 0.0]
						]
					]
				},
				"type": "placeholder1",
				"details": {
					"test": 1234
				}
			}`,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				// No mock setup needed, handler should return error before making db request
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  customErrors.Payload.InvalidPayloadPointError.ErrorMsg,
		},
		{
			name: "Database error on insert",
			inputJSON: `{
				"longlat": {
					"type": "Point",
					"coordinates": [11, 12]
				},
				"type": "placeholder1",
				"details": {
					"test": 1234
				}
			}`,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery("INSERT INTO points").
					WithArgs(
						pgxmock.AnyArg(),
						db.PointType("placeholder1"),
						[]byte(`{"test":1234}`),
					).
					// Simulate a database error
					WillReturnError(fmt.Errorf("Unable to connect to database"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  customErrors.Database.UnknownDatabaseError.ErrorMsg,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runHandlerTest(t, tt, handler.HandlePost, mock)
		})
	}
}
