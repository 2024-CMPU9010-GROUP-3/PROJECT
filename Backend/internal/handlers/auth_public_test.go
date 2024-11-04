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
	userRoute       = "/auth/User"
	userIdString    = "41692803-0f09-4d6b-9b0f-f893bb985bff"
	userIdStringAlt = "48F06E4A-2DF7-4681-B72B-19C99E285CA0"
	username        = `testy`
	usernameAlt     = `testy1`
	email           = `testy@example.com`
	emailAlt        = `testy1@example.com`
	pw              = `test`
	pwLong          = `7z6PX6a1aQieKp94NcaDLCTPBEG1fc90YFZLz4a5rf7TFKMuVEA9trFbgtkpLQUrAEuJp3ffx` // 73 bytes
	firstname       = `Testy`
	firstnameAlt    = `Testy1`
	lastname        = `McTesterson`
	lastnameAlt     = `McTesterson1`
	pwHash          = `$2a$12$oMO4XyesvVS29xYsd8HKn.KBB8J2pxCydSPkuFcTnEfwQaKb2MX2i`
	pfpLink         = `https://www.example.com/image.png`
	pfpLinkAlt      = `https://www.example.com/image1.png`

	queryGetLoginById = `SELECT Id, Username, Email, PasswordHash ` +
		`FROM logins ` +
		`WHERE Id = \$1 ` +
		`LIMIT 1`

	queryGetLoginByEmail = `SELECT Id, Username, Email, PasswordHash ` +
		`FROM logins ` +
		`WHERE Email = \$1 ` +
		`LIMIT 1`

	queryGetLoginByUsername = `SELECT Id, Username, Email, PasswordHash ` +
		`FROM logins ` +
		`WHERE Username = \$1 ` +
		`LIMIT 1`

	queryGetUserDetailsById = `SELECT Id, RegisterDate, FirstName, LastName, ProfilePicture, LastLoggedIn ` +
		`FROM user_details ` +
		`WHERE Id = \$1 ` +
		`LIMIT 1`

	queryInsertIntoLogins = `INSERT INTO logins`

	queryInsertIntoUserDetails = `INSERT INTO user_details`

	queryUpdateLogins = `UPDATE logins`

	queryUpdateUserDetails = `UPDATE user_details`

	jsonCreateUser = `{
				"Username": "%s",
				"Email": "%s",
				"Password": "%s",
				"FirstName": "%s",
				"LastName": "%s",
				"ProfilePicture": "%s"
			}`

	jsonCreateUserUsernameMissing = `{
				"Email": "%s",
				"Password": "%s",
				"FirstName": "%s",
				"LastName": "%s",
				"ProfilePicture": "%s"
			}`

	jsonCreateUserEmailMissing = `{
				"Username": "%s",
				"Password": "%s",
				"FirstName": "%s",
				"LastName": "%s",
				"ProfilePicture": "%s"
			}`

	jsonCreateUserPasswordMissing = `{
				"Email": "%s",
				"Username": "%s",
				"FirstName": "%s",
				"LastName": "%s",
				"ProfilePicture": "%s"
			}`

	jsonCreateUserFirstNameMissing = `{
				"Email": "%s",
				"Username": "%s",
				"Password": "%s",
				"LastName": "%s",
				"ProfilePicture": "%s"
			}`

	jsonCreateUserLastNameMissing = `{
				"Email": "%s",
				"Username": "%s",
				"Password": "%s",
				"FirstName": "%s",
				"ProfilePicture": "%s"
			}`

	jsonCreateUserProfilePictureMissing = `{
				"Email": "%s",
				"Username": "%s",
				"Password": "%s",
				"FirstName": "%s",
				"LastName": "%s"
			}`

	jsonCreateUserFieldOrder = `{
				"FirstName": "%s",
				"Password": "%s",
				"ProfilePicture": "%s",
				"Username": "%s",
				"Email": "%s",
				"LastName": "%s"
			}`

	jsonCreateUserInvalid = `{
				"Username": "%s",
				"Email": "%s",
				"Password": "%s",
				"FirstName": "%s",
				"LastName": "%s",
				"ProfilePicture": "%s"
			`

	jsonResponseUserId = `{
				"error": null,
				"response": {
					"content": {
						"userid": "%s"
					}
				}
			}`

	jsonSimulatedDbError = `{
					"error": {
						"errorCode": 1104,
						"errorMsg": "Unknown database error",
						"cause": "Simulate Database Error"
					},
					"response": null
				}`
)

var (
	rowsGetUserDetails = pgxmock.NewRows([]string{"Id", "RegisterDate", "FirstName", "LastName", "ProfilePicture", "LastLoggedIn"})
	rowsGetLogin       = pgxmock.NewRows([]string{"Id", "Username", "Email", "PasswordHash"})
	rowsId             = pgxmock.NewRows([]string{"Id"})
	simulatedDbError   = fmt.Errorf("Simulate Database Error")
	resultUpdated      = pgconn.NewCommandTag("UPDATED")
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
				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(rowsGetUserDetails.
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
				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(rowsGetUserDetails)
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
				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).WillReturnError(simulatedDbError)
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedJSON:   jsonSimulatedDbError,
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
			Name:      "Positive testcase",
			Method:    "POST",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUser, username, email, pw, firstname, lastname, pfpLink),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnRows(rowsGetLogin)

				mock.ExpectBegin()

				mock.ExpectQuery(queryInsertIntoLogins).
					WithArgs(username, email, testutil.BcryptArg(pw)).
					WillReturnRows(rowsId.
						AddRow(userId))

				mock.ExpectQuery(queryInsertIntoUserDetails).
					WithArgs(userId, firstname, lastname, pfpLinkPG).
					WillReturnRows(rowsId.
						AddRow(userId))

				mock.ExpectCommit()
			},
			ExpectedStatus: http.StatusCreated,
			ExpectedJSON:   fmt.Sprintf(jsonResponseUserId, userIdString),
		},
		{
			Name:      "Json field order should not matter",
			Method:    "POST",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUserFieldOrder, firstname, pw, pfpLink, username, email, lastname),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnRows(rowsGetLogin)

				mock.ExpectBegin()

				mock.ExpectQuery(queryInsertIntoLogins).
					WithArgs(username, email, testutil.BcryptArg(pw)).
					WillReturnRows(rowsId.
						AddRow(userId))

				mock.ExpectQuery(queryInsertIntoUserDetails).
					WithArgs(userId, firstname, lastname, pfpLinkPG).
					WillReturnRows(rowsId.
						AddRow(userId))

				mock.ExpectCommit()
			},
			ExpectedStatus: http.StatusCreated,
			ExpectedJSON:   fmt.Sprintf(jsonResponseUserId, userIdString),
		},
		{
			Name:      "Invalid input json",
			Method:    "POST",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUserInvalid, username, email, pw, firstname, lastname, pfpLink), // json invalid: closing brace missing
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
			Name:      "Username missing",
			Method:    "POST",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUserUsernameMissing, email, pw, firstname, lastname, pfpLink),
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
			Name:      "Email missing",
			Method:    "POST",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUserEmailMissing, username, pw, firstname, lastname, pfpLink),
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
			Name:      "Password missing",
			Method:    "POST",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUserPasswordMissing, email, username, firstname, lastname, pfpLink),
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
			Name:      "Email already exists",
			Method:    "POST",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUser, username, email, pw, firstname, lastname, pfpLink),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(rowsGetLogin.
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
			Name:      "Username already exists",
			Method:    "POST",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUser, username, email, pw, firstname, lastname, pfpLink),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(rowsGetLogin) // no rows returned
				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnRows(rowsGetLogin.
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
			Name:      "Password hashing error",
			Method:    "POST",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUser, username, email, pwLong, firstname, lastname, pfpLink),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(rowsGetLogin) // no rows returned
				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnRows(rowsGetLogin) // no rows returned
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedError:  errors.Payload.PasswordTooLongError.ErrorMsg,
			ExpectedJSON: `{
					"error": {
						"errorCode": 1214,
						"errorMsg": "Password too long (max. 72 bytes)"
					},
					"response": null
				}`,
		},
		{
			Name:      "Begin transaction error",
			Method:    "POST",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUser, username, email, pw, firstname, lastname, pfpLink),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(rowsGetLogin) // no rows returned
				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnRows(rowsGetLogin) // no rows returned

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
			Name:      "Database error during insert",
			Method:    "POST",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUser, username, email, pw, firstname, lastname, pfpLink),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(rowsGetLogin) // no rows returned
				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnRows(rowsGetLogin) // no rows returned

				mock.ExpectBegin()

				mock.ExpectQuery(queryInsertIntoLogins).
					WithArgs(username, email, testutil.BcryptArg(pw)).
					WillReturnError(simulatedDbError)

				mock.ExpectRollback() // transaction should roll bock on error
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedError:  errors.Database.UnknownDatabaseError.ErrorMsg,
			ExpectedJSON:   jsonSimulatedDbError,
		},
		{
			Name:      "Database error during insert (user details)",
			Method:    "POST",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUser, username, email, pw, firstname, lastname, pfpLink),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(rowsGetLogin) // no rows returned
				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnRows(rowsGetLogin) // no rows returned

				mock.ExpectBegin()

				mock.ExpectQuery(queryInsertIntoLogins).
					WithArgs(username, email, testutil.BcryptArg(pw)).
					WillReturnRows(rowsId.
						AddRow(userId))

				mock.ExpectQuery(queryInsertIntoUserDetails).
					WithArgs(userId, firstname, lastname, pfpLinkPG).
					WillReturnError(simulatedDbError)

				mock.ExpectRollback() // transaction should roll bock on error
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedError:  errors.Database.UnknownDatabaseError.ErrorMsg,
			ExpectedJSON:   jsonSimulatedDbError,
		},
		{
			Name:      "Database error during transaction commit",
			Method:    "POST",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUser, username, email, pw, firstname, lastname, pfpLink),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(rowsGetLogin) // no rows returned
				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnRows(rowsGetLogin) // no rows returned

				mock.ExpectBegin()

				mock.ExpectQuery(queryInsertIntoLogins).
					WithArgs(username, email, testutil.BcryptArg(pw)).
					WillReturnRows(rowsId.
						AddRow(userId))

				mock.ExpectQuery(queryInsertIntoUserDetails).
					WithArgs(userId, firstname, lastname, pfpLinkPG).
					WillReturnRows(rowsId.
						AddRow(userId))

				mock.ExpectCommit().WillReturnError(simulatedDbError)

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

	userIdAlt := pgtype.UUID{}
	userIdAlt.Scan(userIdStringAlt)

	pfpLinkPG := pgtype.Text{}
	err = pfpLinkPG.Scan(pfpLink)
	if err != nil {
		t.Fatal(err)
	}

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

	tests := []testutil.HandlerTestDefinition{
		{
			Name:      "Positive testcase",
			Method:    "PUT",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUser, username, email, pw, firstname, lastname, pfpLink),
			PathParams: map[string]string{
				"id": userIdString,
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(rowsGetLogin.
						AddRow(userId, usernameAlt, emailAlt, pwHash))

				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(rowsGetUserDetails.
						AddRow(userId, registerDate, firstname, lastname, pfpLink, lastLoginDate))

				mock.ExpectBegin()

				mock.ExpectExec(queryUpdateLogins).
					WithArgs(userId, username, email, testutil.BcryptArg(pw)).
					WillReturnResult(resultUpdated)

				mock.ExpectExec(queryUpdateUserDetails).
					WithArgs(userId, firstname, lastname, pfpLink).
					WillReturnResult(resultUpdated)

				mock.ExpectCommit()
			},
			ExpectedStatus: http.StatusAccepted,
			ExpectedJSON:   fmt.Sprintf(jsonResponseUserId, userIdString),
		},
		{
			Name:      "Json field order should not matter",
			Method:    "PUT",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUserFieldOrder, firstname, pw, pfpLink, username, email, lastname),
			PathParams: map[string]string{
				"id": userIdString,
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(rowsGetLogin.
						AddRow(userId, usernameAlt, emailAlt, pwHash))

				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(rowsGetUserDetails.
						AddRow(userId, registerDate, firstname, lastname, pfpLink, lastLoginDate))

				mock.ExpectBegin()

				mock.ExpectExec(queryUpdateLogins).
					WithArgs(userId, username, email, testutil.BcryptArg(pw)).
					WillReturnResult(resultUpdated)

				mock.ExpectExec(queryUpdateUserDetails).
					WithArgs(userId, firstname, lastname, pfpLink).
					WillReturnResult(resultUpdated)

				mock.ExpectCommit()
			},
			ExpectedStatus: http.StatusAccepted,
			ExpectedJSON:   fmt.Sprintf(jsonResponseUserId, userIdString),
		},
		{
			Name:      "No error when username is missing",
			Method:    "PUT",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUserUsernameMissing, email, pw, firstname, lastname, pfpLink),
			PathParams: map[string]string{
				"id": userIdString,
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(rowsGetLogin.
						AddRow(userId, usernameAlt, emailAlt, pwHash))

				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(rowsGetUserDetails.
						AddRow(userId, registerDate, firstname, lastname, pfpLink, lastLoginDate))

				mock.ExpectBegin()

				mock.ExpectExec(queryUpdateLogins).
					WithArgs(userId, usernameAlt, email, testutil.BcryptArg(pw)).
					WillReturnResult(resultUpdated)

				mock.ExpectExec(queryUpdateUserDetails).
					WithArgs(userId, firstname, lastname, pfpLink).
					WillReturnResult(resultUpdated)

				mock.ExpectCommit()
			},
			ExpectedStatus: http.StatusAccepted,
			ExpectedJSON:   fmt.Sprintf(jsonResponseUserId, userIdString),
		},
		{
			Name:      "No error when email is missing",
			Method:    "PUT",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUserEmailMissing, username, pw, firstname, lastname, pfpLink),
			PathParams: map[string]string{
				"id": userIdString,
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(rowsGetLogin.
						AddRow(userId, usernameAlt, emailAlt, pwHash))

				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(rowsGetUserDetails.
						AddRow(userId, registerDate, firstname, lastname, pfpLink, lastLoginDate))

				mock.ExpectBegin()

				mock.ExpectExec(queryUpdateLogins).
					WithArgs(userId, username, emailAlt, testutil.BcryptArg(pw)).
					WillReturnResult(resultUpdated)

				mock.ExpectExec(queryUpdateUserDetails).
					WithArgs(userId, firstname, lastname, pfpLink).
					WillReturnResult(resultUpdated)

				mock.ExpectCommit()
			},
			ExpectedStatus: http.StatusAccepted,
			ExpectedJSON:   fmt.Sprintf(jsonResponseUserId, userIdString),
		},
		{
			Name:      "No error when password is missing",
			Method:    "PUT",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUserPasswordMissing, email, username, firstname, lastname, pfpLink),
			PathParams: map[string]string{
				"id": userIdString,
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(rowsGetLogin.
						AddRow(userId, usernameAlt, emailAlt, pwHash))

				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(rowsGetUserDetails.
						AddRow(userId, registerDate, firstname, lastname, pfpLink, lastLoginDate))

				mock.ExpectBegin()

				mock.ExpectExec(queryUpdateLogins).
					WithArgs(userId, username, email, pwHash).
					WillReturnResult(resultUpdated)

				mock.ExpectExec(queryUpdateUserDetails).
					WithArgs(userId, firstname, lastname, pfpLink).
					WillReturnResult(resultUpdated)

				mock.ExpectCommit()
			},
			ExpectedStatus: http.StatusAccepted,
			ExpectedJSON:   fmt.Sprintf(jsonResponseUserId, userIdString),
		},
		{
			Name:      "No error when first name is missing",
			Method:    "PUT",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUserFirstNameMissing, email, username, pw, lastname, pfpLink),
			PathParams: map[string]string{
				"id": userIdString,
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {

				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(rowsGetLogin.
						AddRow(userId, usernameAlt, emailAlt, pwHash))

				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(rowsGetUserDetails.
						AddRow(userId, registerDate, firstnameAlt, lastnameAlt, pfpLinkAlt, lastLoginDate))

				mock.ExpectBegin()

				mock.ExpectExec(queryUpdateLogins).
					WithArgs(userId, username, email, pwHash).
					WillReturnResult(resultUpdated)

				mock.ExpectExec(queryUpdateUserDetails).
					WithArgs(userId, firstnameAlt, lastname, pfpLink).
					WillReturnResult(resultUpdated)

				mock.ExpectCommit()
			},
			ExpectedStatus: http.StatusAccepted,
			ExpectedJSON:   fmt.Sprintf(jsonResponseUserId, userIdString),
		},
		{
			Name:      "No error when last name is missing",
			Method:    "PUT",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUserLastNameMissing, email, username, pw, firstname, pfpLink),
			PathParams: map[string]string{
				"id": userIdString,
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {

				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(rowsGetLogin.
						AddRow(userId, usernameAlt, emailAlt, pwHash))

				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(rowsGetUserDetails.
						AddRow(userId, registerDate, firstnameAlt, lastnameAlt, pfpLinkAlt, lastLoginDate))

				mock.ExpectBegin()

				mock.ExpectExec(queryUpdateLogins).
					WithArgs(userId, username, email, pwHash).
					WillReturnResult(resultUpdated)

				mock.ExpectExec(queryUpdateUserDetails).
					WithArgs(userId, firstname, lastnameAlt, pfpLink).
					WillReturnResult(resultUpdated)

				mock.ExpectCommit()
			},
			ExpectedStatus: http.StatusAccepted,
			ExpectedJSON:   fmt.Sprintf(jsonResponseUserId, userIdString),
		},
		{
			Name:      "No error when profile picture link is missing",
			Method:    "PUT",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUserProfilePictureMissing, email, username, pw, firstname, lastname),
			PathParams: map[string]string{
				"id": userIdString,
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {

				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(rowsGetLogin.
						AddRow(userId, usernameAlt, emailAlt, pwHash))

				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(rowsGetUserDetails.
						AddRow(userId, registerDate, firstnameAlt, lastnameAlt, pfpLinkAlt, lastLoginDate))

				mock.ExpectBegin()

				mock.ExpectExec(queryUpdateLogins).
					WithArgs(userId, username, email, pwHash).
					WillReturnResult(resultUpdated)

				mock.ExpectExec(queryUpdateUserDetails).
					WithArgs(userId, firstname, lastname, pfpLinkAlt).
					WillReturnResult(resultUpdated)

				mock.ExpectCommit()
			},
			ExpectedStatus: http.StatusAccepted,
			ExpectedJSON:   fmt.Sprintf(jsonResponseUserId, userIdString),
		},
		{
			Name:      "User not found",
			Method:    "PUT",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUser, username, email, pw, firstname, lastname, pfpLink),
			PathParams: map[string]string{
				"id": userIdString,
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(rowsGetLogin)
			},
			ExpectedStatus: http.StatusNotFound,
			ExpectedError:  errors.NotFound.UserNotFoundError.ErrorMsg,
			ExpectedJSON: `{
				"error": {
					"code: 1301,
					"errorMsg": "User not found"
				},
				"response": null
			}`,
		},
		{
			Name:      "Id invalid",
			Method:    "PUT",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUser, username, email, pw, firstname, lastname, pfpLink),
			PathParams: map[string]string{
				"id": userIdString[1:], // first char removed
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// handler should return before db calls are made
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  errors.Parameter.InvalidUUIDError.ErrorMsg,
			ExpectedJSON: `{
				"error": {
					"code: 1202,
					"errorMsg": "Parameter invalid, expected type UUIDv4"
				},
				"response": null
			}`,
		},
		{
			Name:      "Invalid payload",
			Method:    "PUT",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUserInvalid, username, email, pw, firstname, lastname, pfpLink), // closing brace removed
			PathParams: map[string]string{
				"id": userIdString,
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// no database calls expected
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  errors.Payload.InvalidPayloadUserError.ErrorMsg,
			ExpectedJSON: `{
				"error": {
					"code: 1212,
					"errorMsg": "Payload (User) not valid"
				},
				"response": null
			}`,
		},
		{
			Name:      "Password too long",
			Method:    "PUT",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUser, username, email, pwLong, firstname, lastname, pfpLink),
			PathParams: map[string]string{
				"id": userIdString,
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// no database calls expected
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  errors.Payload.PasswordTooLongError.ErrorMsg,
			ExpectedJSON: `{
				"error": {
					"code: 1214,
					"errorMsg": "Password too long (max. 72 bytes)"
				},
				"response": null
			}`,
		},
		{
			Name:      "Username already exists",
			Method:    "PUT",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUser, username, email, pw, firstname, lastname, pfpLink),
			PathParams: map[string]string{
				"id": userIdString,
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnRows(rowsGetLogin.
						AddRow(userIdAlt, username, email, pwHash))

			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  errors.Payload.UsernameAlreadyExistsError.ErrorMsg,
			ExpectedJSON: `{
				"error": {
					"code: 1221,
					"errorMsg": "Username already exists"
				},
				"response": null
			}`,
		},
		{
			Name:      "Username already exists, but is same user",
			Method:    "PUT",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUser, username, email, pw, firstname, lastname, pfpLink),
			PathParams: map[string]string{
				"id": userIdString,
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnRows(rowsGetLogin.
						AddRow(userId, username, email, pwHash))

				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(rowsGetLogin.
						AddRow(userId, username, emailAlt, pwHash))

				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(rowsGetUserDetails.
						AddRow(userId, registerDate, firstname, lastname, pfpLink, lastLoginDate))

				mock.ExpectBegin()

				mock.ExpectExec(queryUpdateLogins).
					WithArgs(userId, username, email, testutil.BcryptArg(pw)).
					WillReturnResult(resultUpdated)

				mock.ExpectExec(queryUpdateUserDetails).
					WithArgs(userId, firstname, lastname, pfpLink).
					WillReturnResult(resultUpdated)

				mock.ExpectCommit()

			},
			ExpectedStatus: http.StatusAccepted,
			ExpectedJSON:   fmt.Sprintf(jsonResponseUserId, userIdString),
		},
		{
			Name:      "Email already exists",
			Method:    "PUT",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUser, username, email, pw, firstname, lastname, pfpLink),
			PathParams: map[string]string{
				"id": userIdString,
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(rowsGetLogin.
						AddRow(userIdAlt, usernameAlt, email, pwHash))

			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedError:  errors.Payload.UsernameAlreadyExistsError.ErrorMsg,
			ExpectedJSON: `{
				"error": {
					"code: 1222,
					"errorMsg": "Email already exists"
				},
				"response": null
			}`,
		},
		{
			Name:      "Email already exists, but is same user",
			Method:    "PUT",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUser, username, email, pw, firstname, lastname, pfpLink),
			PathParams: map[string]string{
				"id": userIdString,
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(rowsGetLogin.
						AddRow(userId, usernameAlt, email, pwHash))

				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(rowsGetLogin.
						AddRow(userId, usernameAlt, email, pwHash))

				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(rowsGetUserDetails.
						AddRow(userId, registerDate, firstname, lastname, pfpLink, lastLoginDate))

				mock.ExpectBegin()

				mock.ExpectExec(queryUpdateLogins).
					WithArgs(userId, username, email, testutil.BcryptArg(pw)).
					WillReturnResult(resultUpdated)

				mock.ExpectExec(queryUpdateUserDetails).
					WithArgs(userId, firstname, lastname, pfpLink).
					WillReturnResult(resultUpdated)

				mock.ExpectCommit()

			},
			ExpectedStatus: http.StatusAccepted,
			ExpectedJSON:   fmt.Sprintf(jsonResponseUserId, userIdString),
		},
		{
			Name:      "Error starting transaction",
			Method:    "PUT",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUser, username, email, pw, firstname, lastname, pfpLink),
			PathParams: map[string]string{
				"id": userIdString,
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(rowsGetLogin.
						AddRow(userId, username, email, pwHash))

				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(rowsGetUserDetails.
						AddRow(userId, registerDate, firstname, lastname, pfpLink, lastLoginDate))

				mock.ExpectBegin().WillReturnError(fmt.Errorf("Simulate database errror"))
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedError:  errors.Database.TransactionStartError.ErrorMsg,
			ExpectedJSON: `{
				"error": {
					"code: 1102,
					"errorMsg": "Could not start database transaction"
				},
				"response": null
			}`,
		},
		{
			Name:      "Error updating login",
			Method:    "PUT",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUser, username, email, pw, firstname, lastname, pfpLink),
			PathParams: map[string]string{
				"id": userIdString,
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(rowsGetLogin.
						AddRow(userId, username, email, pwHash))

				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(rowsGetUserDetails.
						AddRow(userId, registerDate, firstname, lastname, pfpLink, lastLoginDate))

				mock.ExpectBegin()

				mock.ExpectExec(queryUpdateLogins).
					WithArgs(userId, username, email, testutil.BcryptArg(pw)).
					WillReturnError(simulatedDbError)

				mock.ExpectRollback()

			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedError:  errors.Database.TransactionStartError.ErrorMsg,
			ExpectedJSON:   jsonSimulatedDbError,
		},
		{
			Name:      "Error updating user details",
			Method:    "PUT",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUser, username, email, pw, firstname, lastname, pfpLink),
			PathParams: map[string]string{
				"id": userIdString,
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(rowsGetLogin.
						AddRow(userId, username, email, pwHash))

				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(rowsGetUserDetails.
						AddRow(userId, registerDate, firstname, lastname, pfpLink, lastLoginDate))

				mock.ExpectBegin()

				mock.ExpectExec(queryUpdateLogins).
					WithArgs(userId, username, email, testutil.BcryptArg(pw)).
					WillReturnResult(resultUpdated)

				mock.ExpectExec(queryUpdateUserDetails).
					WithArgs(userId, firstname, lastname, pfpLink).
					WillReturnError(simulatedDbError)

				mock.ExpectRollback()

			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedError:  errors.Database.TransactionStartError.ErrorMsg,
			ExpectedJSON:   jsonSimulatedDbError,
		},
		{
			Name:      "Error committing transaction",
			Method:    "PUT",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUser, username, email, pw, firstname, lastname, pfpLink),
			PathParams: map[string]string{
				"id": userIdString,
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnRows(rowsGetLogin)

				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(rowsGetLogin.
						AddRow(userId, username, email, pwHash))

				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(rowsGetUserDetails.
						AddRow(userId, registerDate, firstname, lastname, pfpLink, lastLoginDate))

				mock.ExpectBegin()

				mock.ExpectExec(queryUpdateLogins).
					WithArgs(userId, username, email, testutil.BcryptArg(pw)).
					WillReturnResult(resultUpdated)

				mock.ExpectExec(queryUpdateUserDetails).
					WithArgs(userId, firstname, lastname, pfpLink).
					WillReturnResult(resultUpdated)

				mock.ExpectCommit().WillReturnError(simulatedDbError)

				mock.ExpectRollback()

			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedError:  errors.Database.TransactionStartError.ErrorMsg,
			ExpectedJSON:   jsonSimulatedDbError,
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
