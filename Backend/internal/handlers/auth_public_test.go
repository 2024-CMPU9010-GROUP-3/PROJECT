//go:build public

package handlers

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/util/testutil"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v4"
)

const userRoute = "/auth/User"
const userIdString = "41692803-0f09-4d6b-9b0f-f893bb985bff"

func TestAuthHandlerHandleGet(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := context.Background()

	RegisterDatabaseConnection(&ctx, mock)

	authHandler := &AuthHandler{}

	registerDate :=  pgtype.Timestamp{}
	err = registerDate.Scan(time.Date(2024, 10, 10, 12, 34, 56, 789000000, time.UTC))
	if err != nil {
		t.Fatal(err)
	}

	lastLoginDate :=  pgtype.Timestamp{}
	err = lastLoginDate.Scan(time.Date(2024, 10, 30, 12, 34, 56, 789000000, time.UTC))
	if err != nil {
		t.Fatal(err)
	}

	userId := pgtype.UUID{}
	userId.Scan(userIdString)

	tests := []testutil.HandlerTestDefinition {
		{
			Name: "Positive testcase",
			Method: "GET",
			Route: userRoute,
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT Id, RegisterDate, FirstName, LastName, ProfilePicture, LastLoggedIn ` +
												 `FROM user_details ` +
												 `WHERE Id = \$1 ` +
												 `LIMIT 1`).
						 WithArgs(userId).
						 WillReturnRows(pgxmock.NewRows([]string{"Id", "RegisterDate", "FirstName", "LastName", "ProfilePicture", "LastLoggedIn"}).
						 AddRow(userId, registerDate, "Testy", "McTesterson", "https://example.com/", lastLoginDate))
			},
			ExpectedStatus: http.StatusOK,
			ExpectedJSON: fmt.Sprintf(`{
				"error": null,
				"response": {
					"content": {
						"id": "%s",
						"registerdate": "2024-10-10T12:34:56.789Z",
						"firstname": "Testy",
						"lastname": "McTesterson",
						"profilepicture": "https://example.com/",
						"lastloggedin": "2024-10-30T12:34:56.789Z"
					}
				}
			}`, userIdString),
			PathParams: map[string]string{
				"id": userIdString,
			},
		},
		{
			Name: "Invalid UUID",
			Method: "GET",
			Route: userRoute,
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// handler should return before db is queried
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedJSON: `{
					"error": {
						"errorCode": 1202,
						"errorMsg": "Parameter invalid, expected type UUIDv4"
					},
					"response": null
				}`,
			PathParams: map[string]string{
				"id": userIdString[1:], // first char removed
			},
		},
		{
			Name: "User not found",
			Method: "GET",
			Route: userRoute,
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT Id, RegisterDate, FirstName, LastName, ProfilePicture, LastLoggedIn ` +
												 `FROM user_details ` +
												 `WHERE Id = \$1 ` +
												 `LIMIT 1`).
						 WithArgs(userId).
						 WillReturnRows(pgxmock.NewRows([]string{"Id", "RegisterDate", "FirstName", "LastName", "ProfilePicture", "LastLoggedIn"}))
			},
			ExpectedStatus: http.StatusNotFound,
			ExpectedJSON: `{
					"error": {
						"errorCode": 1301,
						"errorMsg": "User not found"
					},
					"response": null
				}`,
			PathParams: map[string]string{
				"id": userIdString,
			},
		},
		{
			Name: "Database error during query",
			Method: "GET",
			Route: userRoute,
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`SELECT Id, RegisterDate, FirstName, LastName, ProfilePicture, LastLoggedIn ` +
												 `FROM user_details ` +
												 `WHERE Id = \$1 ` +
												 `LIMIT 1`).
						 WithArgs(userId).WillReturnError(fmt.Errorf("Simulate Database Error"))
						},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedJSON: `{
					"error": {
						"errorCode": 1104,
						"errorMsg": "Unknown database error",
							"cause": "Simulate Database Error"
					},
					"response": null
				}`,
			PathParams: map[string]string{
				"id": userIdString,
			},
		},
	}
	testutil.RunTests(t, authHandler.HandleGet, mock, tests)
}

func TestAuthHandlerHandlePost(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := context.Background()

	RegisterDatabaseConnection(&ctx, mock)

	authHandler := &AuthHandler{}
	tests := []testutil.HandlerTestDefinition {
	}
	testutil.RunTests(t, authHandler.HandlePost, mock, tests)
}

func TestAuthHandlerHandlePut(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := context.Background()

	RegisterDatabaseConnection(&ctx, mock)
	
	authHandler := &AuthHandler{}
	tests := []testutil.HandlerTestDefinition {
		// TODO: Add test cases
	}
	testutil.RunTests(t, authHandler.HandlePut, mock, tests)
}

func TestAuthHandlerHandleDelete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := context.Background()

	RegisterDatabaseConnection(&ctx, mock)

	authHandler := &AuthHandler{}
	tests := []testutil.HandlerTestDefinition {
		// TODO: Add test cases
	}
	testutil.RunTests(t, authHandler.HandleDelete, mock, tests)
}

func TestAuthHandlerHandleLogin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := context.Background()

	RegisterDatabaseConnection(&ctx, mock)

	authHandler := &AuthHandler{}
	tests := []testutil.HandlerTestDefinition {
		// TODO: Add test cases
	}
	testutil.RunTests(t, authHandler.HandleLogin, mock, tests)
}
