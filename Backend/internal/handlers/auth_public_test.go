//go:build public

package handlers

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/errors"
	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/util/testutil"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v4"
)

const (
	userRoute    = "/auth/User"
	userIdString = "41692803-0f09-4d6b-9b0f-f893bb985bff"
	username     = `testy`
	email        = `testy@example.com`
	pw           = `test`
	pwLong       = `7z6PX6a1aQieKp94NcaDLCTPBEG1fc90YFZLz4a5rf7TFKMuVEA9trFbgtkpLQUrAEuJp3ffx` // 73 bytes
	firstname    = `Testy`
	lastname     = `McTesterson`
	pwHash       = `$2a$12$oMO4XyesvVS29xYsd8HKn.KBB8J2pxCydSPkuFcTnEfwQaKb2MX2i`
	pfpLink      = `https://www.example.com/image.png`
)

func TestAuthHandlerHandleGet(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	ctx := context.Background()

	RegisterDatabaseConnection(&ctx, mock)

	authHandler := &AuthHandler{}

	registerDate := pgtype.Timestamp{}
	err = registerDate.Scan(time.Date(2024, 10, 10, 12, 34, 56, 789000000, time.UTC))
	if err != nil {
		t.Fatal(err)
	}

	lastLoginDate := pgtype.Timestamp{}
	err = lastLoginDate.Scan(time.Date(2024, 10, 30, 12, 34, 56, 789000000, time.UTC))
	if err != nil {
		t.Fatal(err)
	}

	userId := pgtype.UUID{}
	userId.Scan(userIdString)

	tests := []testutil.HandlerTestDefinition{
		{
			Name:   "Positive testcase",
			Method: "GET",
			Route:  userRoute,
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
			Name:   "Invalid UUID",
			Method: "GET",
			Route:  userRoute,
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
			Name:   "User not found",
			Method: "GET",
			Route:  userRoute,
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
			Name:   "Database error during query",
			Method: "GET",
			Route:  userRoute,
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

	userId := pgtype.UUID{}
	userId.Scan(userIdString)

	pfpLinkPG := pgtype.Text{}
	err = pfpLinkPG.Scan(pfpLink)
	if err != nil {
		t.Fatal(err)
	}

	tests := []testutil.HandlerTestDefinition{
		{
			Name:   "Positive testcase",
			Method: "POST",
			Route:  userRoute,
			InputJSON: fmt.Sprintf(`{
				"Username": "%s",
				"Email": "%s",
				"Password": "%s",
				"FirstName": "%s",
				"LastName": "%s",
				"ProfilePicture": "%s"
			}`, username, email, pw, firstname, lastname, pfpLink),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(
					`SELECT Id, Username, Email, PasswordHash ` +
						`FROM logins ` +
						`WHERE Email = \$1 ` +
						`LIMIT 1`).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Username", "Email", "PasswordHash"}))

				mock.ExpectQuery(
					`SELECT Id, Username, Email, PasswordHash ` +
						`FROM logins ` +
						`WHERE Username = \$1 ` +
						`LIMIT 1`).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Username", "Email", "PasswordHash"}))

				mock.ExpectBegin()

				mock.ExpectQuery(`INSERT INTO logins`).
					WithArgs(username, email, testutil.BcryptArg(pw)).
					WillReturnRows(pgxmock.NewRows([]string{"Id"}).
						AddRow(userId))

				mock.ExpectQuery(`INSERT INTO user_details`).
					WithArgs(userId, firstname, lastname, pfpLinkPG).
					WillReturnRows(pgxmock.NewRows([]string{"Id"}).
						AddRow(userId))

				mock.ExpectCommit()
			},
			ExpectedStatus: http.StatusCreated,
			ExpectedJSON: fmt.Sprintf(`{
				"error": null,
				"response": {
					"content": {
						"userid": "%s"
					}
				}
			}`, userIdString),
		},
		{
			Name:   "Json field order should not matter",
			Method: "POST",
			Route:  userRoute,
			InputJSON: fmt.Sprintf(`{
				"FirstName": "%s",
				"Password": "%s",
				"ProfilePicture": "%s",
				"Username": "%s",
				"Email": "%s",
				"LastName": "%s"
			}`, firstname, pw, pfpLink, username, email, lastname),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(
					`SELECT Id, Username, Email, PasswordHash ` +
						`FROM logins ` +
						`WHERE Email = \$1 ` +
						`LIMIT 1`).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Username", "Email", "PasswordHash"}))

				mock.ExpectQuery(
					`SELECT Id, Username, Email, PasswordHash ` +
						`FROM logins ` +
						`WHERE Username = \$1 ` +
						`LIMIT 1`).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Username", "Email", "PasswordHash"}))

				mock.ExpectBegin()

				mock.ExpectQuery(`INSERT INTO logins`).
					WithArgs(username, email, testutil.BcryptArg(pw)).
					WillReturnRows(pgxmock.NewRows([]string{"Id"}).
						AddRow(userId))

				mock.ExpectQuery(`INSERT INTO user_details`).
					WithArgs(userId, firstname, lastname, pfpLinkPG).
					WillReturnRows(pgxmock.NewRows([]string{"Id"}).
						AddRow(userId))

				mock.ExpectCommit()
			},
			ExpectedStatus: http.StatusCreated,
			ExpectedJSON: fmt.Sprintf(`{
				"error": null,
				"response": {
					"content": {
						"userid": "%s"
					}
				}
			}`, userIdString),
		},
		{
			Name:   "Invalid input json",
			Method: "POST",
			Route:  userRoute,
			InputJSON: fmt.Sprintf(`{
				"Username": "%s",
				"Email": "%s",
				"Password": "%s",
				"FirstName": "%s",
				"LastName": "%s",
				"ProfilePicture": "%s"
			`, username, email, pw, firstname, lastname, pfpLink), // json invalid: closing brace missing
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// Handler should return error before db call is made
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  errors.Payload.InvalidPayloadUserError.ErrorMsg,
			ExpectedJSON: `{
					"error": {
						"errorCode": 1212,
						"errorMsg": "Payload (User) not valid"
					},
					"response": null
				}`,
		},
		{
			Name:   "Username missing",
			Method: "POST",
			Route:  userRoute,
			InputJSON: fmt.Sprintf(`{
				"Email": "%s",
				"Password": "%s",
				"FirstName": "%s",
				"LastName": "%s",
				"ProfilePicture": "%s"
			}`, email, pw, firstname, lastname, pfpLink),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// Handler should return error before db call is made
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  errors.Parameter.RequiredParameterMissingError.ErrorMsg,
			ExpectedJSON: `{
					"error": {
						"errorCode": 1201,
						"errorMsg": "One or more required parameters are missing",
						"cause":"Username is required"
					},
					"response": null
				}`,
		},
		{
			Name:   "Email missing",
			Method: "POST",
			Route:  userRoute,
			InputJSON: fmt.Sprintf(`{
				"Username": "%s",
				"Password": "%s",
				"FirstName": "%s",
				"LastName": "%s",
				"ProfilePicture": "%s"
			}`, username, pw, firstname, lastname, pfpLink),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// Handler should return error before db call is made
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  errors.Parameter.RequiredParameterMissingError.ErrorMsg,
			ExpectedJSON: `{
					"error": {
						"errorCode": 1201,
						"errorMsg": "One or more required parameters are missing",
						"cause":"Email is required"
					},
					"response": null
				}`,
		},
		{
			Name:   "Password missing",
			Method: "POST",
			Route:  userRoute,
			InputJSON: fmt.Sprintf(`{
				"Email": "%s",
				"Username": "%s",
				"FirstName": "%s",
				"LastName": "%s",
				"ProfilePicture": "%s"
			}`, email, username, firstname, lastname, pfpLink),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// Handler should return error before db call is made
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  errors.Parameter.RequiredParameterMissingError.ErrorMsg,
			ExpectedJSON: `{
					"error": {
						"errorCode": 1201,
						"errorMsg": "One or more required parameters are missing",
						"cause":"Password is required"
					},
					"response": null
				}`,
		},
		{
			Name:   "Email already exists",
			Method: "POST",
			Route:  userRoute,
			InputJSON: fmt.Sprintf(`{
				"Email": "%s",
				"Username": "%s",
				"Password": "%s",
				"FirstName": "%s",
				"LastName": "%s",
				"ProfilePicture": "%s"
			}`, email, username, pw, firstname, lastname, pfpLink),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(
					`SELECT Id, Username, Email, PasswordHash ` +
						`FROM logins ` +
						`WHERE Email = \$1 ` +
						`LIMIT 1`).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Username", "Email", "PasswordHash"}).
						AddRow(userId, username, email, pwHash))
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  errors.Payload.EmailAlreadyExistsError.ErrorMsg,
			ExpectedJSON: `{
					"error": {
						"errorCode": 1222,
						"errorMsg": "Email already exists"
					},
					"response": null
				}`,
		},
		{
			Name:   "Username already exists",
			Method: "POST",
			Route:  userRoute,
			InputJSON: fmt.Sprintf(`{
				"Email": "%s",
				"Username": "%s",
				"Password": "%s",
				"FirstName": "%s",
				"LastName": "%s",
				"ProfilePicture": "%s"
			}`, email, username, pw, firstname, lastname, pfpLink),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(
					`SELECT Id, Username, Email, PasswordHash ` +
						`FROM logins ` +
						`WHERE Email = \$1 ` +
						`LIMIT 1`).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Username", "Email", "PasswordHash"})) // no rows returned
				mock.ExpectQuery(
					`SELECT Id, Username, Email, PasswordHash ` +
						`FROM logins ` +
						`WHERE Username = \$1 ` +
						`LIMIT 1`).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Username", "Email", "PasswordHash"}).
						AddRow(userId, username, email, pwHash))
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  errors.Payload.UsernameAlreadyExistsError.ErrorMsg,
			ExpectedJSON: `{
					"error": {
						"errorCode": 1221,
						"errorMsg": "Username already exists"
					},
					"response": null
				}`,
		},
		{
			Name:   "Password hashing error",
			Method: "POST",
			Route:  userRoute,
			InputJSON: fmt.Sprintf(`{
				"Email": "%s",
				"Username": "%s",
				"Password": "%s",
				"FirstName": "%s",
				"LastName": "%s",
				"ProfilePicture": "%s"
			}`, email, username, pwLong, firstname, lastname, pfpLink),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(
					`SELECT Id, Username, Email, PasswordHash ` +
						`FROM logins ` +
						`WHERE Email = \$1 ` +
						`LIMIT 1`).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Username", "Email", "PasswordHash"})) // no rows returned
				mock.ExpectQuery(
					`SELECT Id, Username, Email, PasswordHash ` +
						`FROM logins ` +
						`WHERE Username = \$1 ` +
						`LIMIT 1`).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Username", "Email", "PasswordHash"})) // no rows returned
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedError:  errors.Internal.HashingError.ErrorMsg,
			ExpectedJSON: `{
					"error": {
						"errorCode": 1012,
						"errorMsg": "Could not hash password"
					},
					"response": null
				}`,
		},
		{
			Name:   "Begin transaction error",
			Method: "POST",
			Route:  userRoute,
			InputJSON: fmt.Sprintf(`{
				"Email": "%s",
				"Username": "%s",
				"Password": "%s",
				"FirstName": "%s",
				"LastName": "%s",
				"ProfilePicture": "%s"
			}`, email, username, pw, firstname, lastname, pfpLink),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(
					`SELECT Id, Username, Email, PasswordHash ` +
						`FROM logins ` +
						`WHERE Email = \$1 ` +
						`LIMIT 1`).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Username", "Email", "PasswordHash"})) // no rows returned
				mock.ExpectQuery(
					`SELECT Id, Username, Email, PasswordHash ` +
						`FROM logins ` +
						`WHERE Username = \$1 ` +
						`LIMIT 1`).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Username", "Email", "PasswordHash"})) // no rows returned

				mock.ExpectBegin().WillReturnError(fmt.Errorf("Could not start database transaction"))
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedError:  errors.Database.TransactionStartError.ErrorMsg,
			ExpectedJSON: `{
					"error": {
						"errorCode": 1102,
						"errorMsg": "Could not start database transaction"
					},
					"response": null
				}`,
		},
		{
			Name:   "Database error during insert",
			Method: "POST",
			Route:  userRoute,
			InputJSON: fmt.Sprintf(`{
				"Email": "%s",
				"Username": "%s",
				"Password": "%s",
				"FirstName": "%s",
				"LastName": "%s",
				"ProfilePicture": "%s"
			}`, email, username, pw, firstname, lastname, pfpLink),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(
					`SELECT Id, Username, Email, PasswordHash ` +
						`FROM logins ` +
						`WHERE Email = \$1 ` +
						`LIMIT 1`).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Username", "Email", "PasswordHash"})) // no rows returned
				mock.ExpectQuery(
					`SELECT Id, Username, Email, PasswordHash ` +
						`FROM logins ` +
						`WHERE Username = \$1 ` +
						`LIMIT 1`).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Username", "Email", "PasswordHash"})) // no rows returned

				mock.ExpectBegin()

				mock.ExpectQuery(`INSERT INTO logins`).
					WithArgs(username, email, testutil.BcryptArg(pw)).
					WillReturnError(fmt.Errorf("Simulate database error"))

				mock.ExpectRollback() // transaction should roll bock on error
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedError:  errors.Database.UnknownDatabaseError.ErrorMsg,
			ExpectedJSON: `{
					"error": {
						"errorCode": 1104,
						"errorMsg": "Unknown database error",
            "cause": "Simulate database error"
					},
					"response": null
				}`,
		},
		{
			Name:   "Database error during insert (user details)",
			Method: "POST",
			Route:  userRoute,
			InputJSON: fmt.Sprintf(`{
				"Email": "%s",
				"Username": "%s",
				"Password": "%s",
				"FirstName": "%s",
				"LastName": "%s",
				"ProfilePicture": "%s"
			}`, email, username, pw, firstname, lastname, pfpLink),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(
					`SELECT Id, Username, Email, PasswordHash ` +
						`FROM logins ` +
						`WHERE Email = \$1 ` +
						`LIMIT 1`).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Username", "Email", "PasswordHash"})) // no rows returned
				mock.ExpectQuery(
					`SELECT Id, Username, Email, PasswordHash ` +
						`FROM logins ` +
						`WHERE Username = \$1 ` +
						`LIMIT 1`).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Username", "Email", "PasswordHash"})) // no rows returned

				mock.ExpectBegin()

				mock.ExpectQuery(`INSERT INTO logins`).
					WithArgs(username, email, testutil.BcryptArg(pw)).
					WillReturnRows(pgxmock.NewRows([]string{"Id"}).
						AddRow(userId))

				mock.ExpectQuery(`INSERT INTO user_details`).
					WithArgs(userId, firstname, lastname, pfpLinkPG).
					WillReturnError(fmt.Errorf("Simulate database error"))

				mock.ExpectRollback() // transaction should roll bock on error
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedError:  errors.Database.UnknownDatabaseError.ErrorMsg,
			ExpectedJSON: `{
					"error": {
						"errorCode": 1104,
						"errorMsg": "Unknown database error",
            "cause": "Simulate database error"
					},
					"response": null
				}`,
		},
		{
			Name:   "Database error during transaction commit",
			Method: "POST",
			Route:  userRoute,
			InputJSON: fmt.Sprintf(`{
				"Email": "%s",
				"Username": "%s",
				"Password": "%s",
				"FirstName": "%s",
				"LastName": "%s",
				"ProfilePicture": "%s"
			}`, email, username, pw, firstname, lastname, pfpLink),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(
					`SELECT Id, Username, Email, PasswordHash ` +
						`FROM logins ` +
						`WHERE Email = \$1 ` +
						`LIMIT 1`).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Username", "Email", "PasswordHash"})) // no rows returned
				mock.ExpectQuery(
					`SELECT Id, Username, Email, PasswordHash ` +
						`FROM logins ` +
						`WHERE Username = \$1 ` +
						`LIMIT 1`).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Username", "Email", "PasswordHash"})) // no rows returned

				mock.ExpectBegin()

				mock.ExpectQuery(`INSERT INTO logins`).
					WithArgs(username, email, testutil.BcryptArg(pw)).
					WillReturnRows(pgxmock.NewRows([]string{"Id"}).
						AddRow(userId))

				mock.ExpectQuery(`INSERT INTO user_details`).
					WithArgs(userId, firstname, lastname, pfpLinkPG).
					WillReturnRows(pgxmock.NewRows([]string{"Id"}).
						AddRow(userId))

				mock.ExpectCommit().WillReturnError(fmt.Errorf("Simulate database error"))

				mock.ExpectRollback() // transaction should roll bock on error
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedError:  errors.Database.TransactionCommitError.ErrorMsg,
			ExpectedJSON: `{
					"error": {
						"errorCode": 1103,
						"errorMsg": "Could not commit database transaction"
					},
					"response": null
				}`,
		},
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

	userId := pgtype.UUID{}
	userId.Scan(userIdString)

	pfpLinkPG := pgtype.Text{}
	err = pfpLinkPG.Scan(pfpLink)
	if err != nil {
		t.Fatal(err)
	}

	tests := []testutil.HandlerTestDefinition{
		{
			Name:   "Positive testcase",
			Method: "POST",
			Route:  userRoute,
			InputJSON: fmt.Sprintf(`{
				"Username": "%s",
				"Email": "%s",
				"Password": "%s",
				"FirstName": "%s",
				"LastName": "%s",
				"ProfilePicture": "%s"
			}`, username, email, pw, firstname, lastname, pfpLink),
			PathParams: map[string]string{
				"id": userIdString,
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(
					`SELECT Id, Username, Email, PasswordHash ` +
						`FROM logins ` +
						`WHERE Email = \$1 ` +
						`LIMIT 1`).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Username", "Email", "PasswordHash"}))

				mock.ExpectQuery(
					`SELECT Id, Username, Email, PasswordHash ` +
						`FROM logins ` +
						`WHERE Username = \$1 ` +
						`LIMIT 1`).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Username", "Email", "PasswordHash"}))

				mock.ExpectBegin()

				mock.ExpectExec(`UPDATE logins`).
					WithArgs(userId, username, email, testutil.BcryptArg(pw)).
					WillReturnResult(pgconn.NewCommandTag("UPDATED"))
					
					mock.ExpectExec(`UPDATE user_details`).
					WithArgs(userId, firstname, lastname, pfpLink).
					WillReturnResult(pgconn.NewCommandTag("UPDATED"))

				mock.ExpectCommit()
			},
			ExpectedStatus: http.StatusAccepted,
			ExpectedJSON: fmt.Sprintf(`{
				"error": null,
				"response": {
					"content": {
						"userid": "%s"
					}
				}
			}`, userIdString),
		},
		{
			Name:   "Json field order should not matter",
			Method: "POST",
			Route:  userRoute,
			InputJSON: fmt.Sprintf(`{
				"FirstName": "%s",
				"Password": "%s",
				"ProfilePicture": "%s",
				"Username": "%s",
				"Email": "%s",
				"LastName": "%s"
			}`, firstname, pw, pfpLink, username, email, lastname),
			PathParams: map[string]string{
				"id": userIdString,
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(
					`SELECT Id, Username, Email, PasswordHash ` +
						`FROM logins ` +
						`WHERE Email = \$1 ` +
						`LIMIT 1`).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Username", "Email", "PasswordHash"}))

				mock.ExpectQuery(
					`SELECT Id, Username, Email, PasswordHash ` +
						`FROM logins ` +
						`WHERE Username = \$1 ` +
						`LIMIT 1`).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Username", "Email", "PasswordHash"}))

				mock.ExpectBegin()

				mock.ExpectExec(`UPDATE logins`).
					WithArgs(userId, username, email, testutil.BcryptArg(pw)).
					WillReturnResult(pgconn.NewCommandTag("UPDATED"))
					
					mock.ExpectExec(`UPDATE user_details`).
					WithArgs(userId, firstname, lastname, pfpLink).
					WillReturnResult(pgconn.NewCommandTag("UPDATED"))

				mock.ExpectCommit()
			},
			ExpectedStatus: http.StatusAccepted,
			ExpectedJSON: fmt.Sprintf(`{
				"error": null,
				"response": {
					"content": {
						"userid": "%s"
					}
				}
			}`, userIdString),
		},
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
	tests := []testutil.HandlerTestDefinition{
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
	tests := []testutil.HandlerTestDefinition{
		// TODO: Add test cases
	}
	testutil.RunTests(t, authHandler.HandleLogin, mock, tests)
}
