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
	loginRoute      = "/auth/User/login"
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
	pwHashAlt       = `$2a$12$iq.hG5JXuvVSPUHSoEuEaOQ4SdFjtzk7PSFzJLjJbOS0j4YgWiPxm`
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

	queryGetEmailExists = `SELECT EXISTS(` +
		`SELECT 1 ` +
		`FROM logins ` +
		`WHERE Email = \$1 ` +
		`AND Id <> \$2) AS "exists"`

	queryGetUsernameExists = `SELECT EXISTS(` +
		`SELECT 1 ` +
		`FROM logins ` +
		`WHERE Username = \$1 ` +
		`AND Id <> \$2) AS "exists"`

	queryGetUserDetailsById = `SELECT Id, RegisterDate, FirstName, LastName, ProfilePicture, LastLoggedIn ` +
		`FROM user_details ` +
		`WHERE Id = \$1 ` +
		`LIMIT 1`

	queryUpdateLastLogin = `UPDATE user_details ` +
		`SET LastLoggedIn = (NOW() AT TIME ZONE 'utc') ` +
		`WHERE Id = \$1`

	queryInsertIntoLogins = `INSERT INTO logins`

	queryInsertIntoUserDetails = `INSERT INTO user_details`

	queryUpdateLogins = `UPDATE logins`

	queryUpdateUserDetails = `UPDATE user_details`

	queryDeleteFromUsers = `DELETE FROM logins WHERE Id = \$1 CASCADE`

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

	jsonLoginUserWithEmail = `{
				"Email": "%s",
				"Password": "%s"
			}`

	jsonLoginUserWithUsername = `{
				"Username": "%s",
				"Password": "%s"
			}`

	jsonLoginUserWithUsernameInvalid = `{
				"Username": "%s",
				"Password": "%s"
			`

	jsonLoginUserWithUsernamePasswordMissing = `{
				"Username": "%s"
			}`

	jsonLoginUserUsernameAndEmailMissing = `{
				"Password": "%s"
			}`

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

	jsonInvalidUUIDError = `{
					"error": {
						"errorCode": 1202,
						"errorMsg": "Parameter invalid, expected type UUIDv4"
					},
					"response": null
				}`

	jsonWrongCredentialsError = `{
					"error": {
						"errorCode": 1402,
						"errorMsg": "Wrong username or password"
					},
					"response": null
				}`

	jsonPasswordRequiredError = `{
					"error": {
						"errorCode": 1201,
						"errorMsg": "One or more required parameters are missing",
						"cause":"Password is required"
					},
					"response": null
				}`

	jsonUsernameOrEmailRequiredError = `{
					"error": {
						"errorCode": 1201,
						"errorMsg": "One or more required parameters are missing",
						"cause":"Username or Email is required"
					},
					"response": null
				}`

	jsonInvalidUserPayloadError = `{
					"error": {
						"errorCode": 1212,
						"errorMsg": "Payload (User) not valid"
					},
					"response": null
				}`
)

var (
	rowsGetUserDetails = []string{"Id", "RegisterDate", "FirstName", "LastName", "ProfilePicture", "LastLoggedIn"}
	rowsGetLogin       = []string{"Id", "Username", "Email", "PasswordHash"}
	rowsGetExists      = []string{"exists"}
	rowsId             = []string{"Id"}
	simulatedDbError   = fmt.Errorf("Simulate Database Error")
	resultUpdated      = pgconn.NewCommandTag("UPDATED")
	resultDeleted      = pgconn.NewCommandTag("DELETED")
	defaultEnv         = map[string]string{
		"MAGPIE_JWT_SECRET": `RyA4diC7nVdi39Isb9UlujsKN6/qyEjPFVHeLA9VakA=`,
		"MAGPIE_JWT_EXPIRY": `168h`,
	}
	jwtSecretMissingEnv = map[string]string{
		"MAGPIE_JWT_EXPIRY": `168h`,
	}
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
					WillReturnRows(pgxmock.NewRows(rowsGetUserDetails).
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
			ExpectedJSON:   jsonInvalidUUIDError,
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
					WillReturnRows(pgxmock.NewRows(rowsGetUserDetails))
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
				mock.ExpectQuery(queryGetEmailExists).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectQuery(queryGetUsernameExists).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectBegin()

				mock.ExpectQuery(queryInsertIntoLogins).
					WithArgs(username, email, testutil.BcryptArg(pw)).
					WillReturnRows(pgxmock.NewRows(rowsId).
						AddRow(userId))

				mock.ExpectQuery(queryInsertIntoUserDetails).
					WithArgs(userId, firstname, lastname, pfpLinkPG).
					WillReturnRows(pgxmock.NewRows(rowsId).
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
				mock.ExpectQuery(queryGetEmailExists).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectQuery(queryGetUsernameExists).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectBegin()

				mock.ExpectQuery(queryInsertIntoLogins).
					WithArgs(username, email, testutil.BcryptArg(pw)).
					WillReturnRows(pgxmock.NewRows(rowsId).
						AddRow(userId))

				mock.ExpectQuery(queryInsertIntoUserDetails).
					WithArgs(userId, firstname, lastname, pfpLinkPG).
					WillReturnRows(pgxmock.NewRows(rowsId).
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
			ExpectedJSON:   jsonInvalidUserPayloadError,
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
				mock.ExpectQuery(queryGetEmailExists).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(true))
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
				mock.ExpectQuery(queryGetEmailExists).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))
				mock.ExpectQuery(queryGetUsernameExists).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(true))
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
				mock.ExpectQuery(queryGetEmailExists).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))
				mock.ExpectQuery(queryGetUsernameExists).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))
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
				mock.ExpectQuery(queryGetEmailExists).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))
				mock.ExpectQuery(queryGetUsernameExists).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

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
				mock.ExpectQuery(queryGetEmailExists).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))
				mock.ExpectQuery(queryGetUsernameExists).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

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
				mock.ExpectQuery(queryGetEmailExists).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))
				mock.ExpectQuery(queryGetUsernameExists).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectBegin()

				mock.ExpectQuery(queryInsertIntoLogins).
					WithArgs(username, email, testutil.BcryptArg(pw)).
					WillReturnRows(pgxmock.NewRows(rowsId).
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
				mock.ExpectQuery(queryGetEmailExists).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))
				mock.ExpectQuery(queryGetUsernameExists).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectBegin()

				mock.ExpectQuery(queryInsertIntoLogins).
					WithArgs(username, email, testutil.BcryptArg(pw)).
					WillReturnRows(pgxmock.NewRows(rowsId).
						AddRow(userId))

				mock.ExpectQuery(queryInsertIntoUserDetails).
					WithArgs(userId, firstname, lastname, pfpLinkPG).
					WillReturnRows(pgxmock.NewRows(rowsId).
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
				mock.ExpectQuery(queryGetEmailExists).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectQuery(queryGetUsernameExists).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(pgxmock.NewRows(rowsGetLogin).
						AddRow(userId, usernameAlt, emailAlt, pwHash))

				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(pgxmock.NewRows(rowsGetUserDetails).
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
				mock.ExpectQuery(queryGetEmailExists).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectQuery(queryGetUsernameExists).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(pgxmock.NewRows(rowsGetLogin).
						AddRow(userId, usernameAlt, emailAlt, pwHash))

				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(pgxmock.NewRows(rowsGetUserDetails).
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
					WillReturnRows(pgxmock.NewRows(rowsGetLogin).
						AddRow(userId, usernameAlt, emailAlt, pwHash))

				mock.ExpectQuery(queryGetEmailExists).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(pgxmock.NewRows(rowsGetUserDetails).
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
				mock.ExpectQuery(queryGetUsernameExists).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows(rowsGetLogin))

				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(pgxmock.NewRows(rowsGetLogin).
						AddRow(userId, usernameAlt, emailAlt, pwHash))

				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(pgxmock.NewRows(rowsGetUserDetails).
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
				mock.ExpectQuery(queryGetEmailExists).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectQuery(queryGetUsernameExists).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(pgxmock.NewRows(rowsGetLogin).
						AddRow(userId, usernameAlt, emailAlt, pwHash))

				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(pgxmock.NewRows(rowsGetUserDetails).
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

				mock.ExpectQuery(queryGetEmailExists).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectQuery(queryGetUsernameExists).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(pgxmock.NewRows(rowsGetLogin).
						AddRow(userId, usernameAlt, emailAlt, pwHash))

				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(pgxmock.NewRows(rowsGetUserDetails).
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

				mock.ExpectQuery(queryGetEmailExists).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectQuery(queryGetUsernameExists).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(pgxmock.NewRows(rowsGetLogin).
						AddRow(userId, usernameAlt, emailAlt, pwHash))

				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(pgxmock.NewRows(rowsGetUserDetails).
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

				mock.ExpectQuery(queryGetEmailExists).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectQuery(queryGetUsernameExists).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(pgxmock.NewRows(rowsGetLogin).
						AddRow(userId, usernameAlt, emailAlt, pwHash))

				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(pgxmock.NewRows(rowsGetUserDetails).
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
				mock.ExpectQuery(queryGetEmailExists).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectQuery(queryGetUsernameExists).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(pgxmock.NewRows(rowsGetLogin))
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
			ExpectedJSON:   jsonInvalidUUIDError,
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
			ExpectedJSON:   jsonInvalidUserPayloadError,
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
				mock.ExpectQuery(queryGetEmailExists).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectQuery(queryGetUsernameExists).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(true))

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
			Name:      "Email already exists",
			Method:    "PUT",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUser, username, email, pw, firstname, lastname, pfpLink),
			PathParams: map[string]string{
				"id": userIdString,
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetEmailExists).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(true))

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
			Name:      "Error starting transaction",
			Method:    "PUT",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUser, username, email, pw, firstname, lastname, pfpLink),
			PathParams: map[string]string{
				"id": userIdString,
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetEmailExists).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectQuery(queryGetUsernameExists).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(pgxmock.NewRows(rowsGetLogin).
						AddRow(userId, username, email, pwHash))

				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(pgxmock.NewRows(rowsGetUserDetails).
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
				mock.ExpectQuery(queryGetEmailExists).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectQuery(queryGetUsernameExists).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(pgxmock.NewRows(rowsGetLogin).
						AddRow(userId, username, email, pwHash))

				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(pgxmock.NewRows(rowsGetUserDetails).
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
				mock.ExpectQuery(queryGetEmailExists).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectQuery(queryGetUsernameExists).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(pgxmock.NewRows(rowsGetLogin).
						AddRow(userId, username, email, pwHash))

				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(pgxmock.NewRows(rowsGetUserDetails).
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
				mock.ExpectQuery(queryGetEmailExists).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectQuery(queryGetUsernameExists).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows(rowsGetExists).AddRow(false))

				mock.ExpectQuery(queryGetLoginById).
					WithArgs(userId).
					WillReturnRows(pgxmock.NewRows(rowsGetLogin).
						AddRow(userId, username, email, pwHash))

				mock.ExpectQuery(queryGetUserDetailsById).
					WithArgs(userId).
					WillReturnRows(pgxmock.NewRows(rowsGetUserDetails).
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

	userId := pgtype.UUID{}
	userId.Scan(userIdString)

	tests := []testutil.HandlerTestDefinition{
		{
			Name:      "Positive testcase",
			Method:    "DELETE",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUser, username, email, pw, firstname, lastname, pfpLink),
			PathParams: map[string]string{
				"id": userIdString,
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(queryDeleteFromUsers).WithArgs(userId).WillReturnResult(resultDeleted)
			},
			ExpectedStatus: http.StatusAccepted,
			ExpectedJSON:   fmt.Sprintf(jsonResponseUserId, userIdString),
		},
		{
			Name:      "Invalid UUID",
			Method:    "DELETE",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUser, username, email, pw, firstname, lastname, pfpLink),
			PathParams: map[string]string{
				"id": userIdString[1:],
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// should return before any calls are made
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedJSON:   jsonInvalidUUIDError,
		},
		{
			Name:      "Error during delete",
			Method:    "DELETE",
			Route:     userRoute,
			InputJSON: fmt.Sprintf(jsonCreateUser, username, email, pw, firstname, lastname, pfpLink),
			PathParams: map[string]string{
				"id": userIdString[1:],
			},
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(queryDeleteFromUsers).WithArgs(userId).WillReturnError(simulatedDbError)
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedJSON:   jsonSimulatedDbError,
		},
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

	userId := pgtype.UUID{}
	userId.Scan(userIdString)

	duration, err := time.ParseDuration(`168h`)
	if err != nil {
		t.Errorf("Could not parse jwt duration: %v", err)
	}

	tests := []testutil.HandlerTestDefinition{
		{
			Name:      "Positive testcase (email)",
			Method:    "POST",
			Route:     loginRoute,
			Env:       defaultEnv,
			InputJSON: fmt.Sprintf(jsonLoginUserWithEmail, email, pw),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetLogin).AddRow(userId, username, email, pwHash))

				mock.ExpectExec(queryUpdateLastLogin).WithArgs(userId).WillReturnResult(resultUpdated)
			},
			ExpectedCookies: []*http.Cookie{
				{
					Name:     "magpie_auth",
					HttpOnly: true,
					SameSite: http.SameSiteLaxMode,
					Expires:  time.Now().Add(duration),
					Path:     "/",
				},
			},
			ExpectedStatus: http.StatusOK,
			ExpectedJSON:   fmt.Sprintf(jsonResponseUserId, userIdString),
		},
		{
			Name:      "Positive testcase (username)",
			Method:    "POST",
			Route:     loginRoute,
			Env:       defaultEnv,
			InputJSON: fmt.Sprintf(jsonLoginUserWithUsername, username, pw),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows(rowsGetLogin).AddRow(userId, username, email, pwHash))

				mock.ExpectExec(queryUpdateLastLogin).WithArgs(userId).WillReturnResult(resultUpdated)
			},
			ExpectedCookies: []*http.Cookie{
				{
					Name:     "magpie_auth",
					HttpOnly: true,
					SameSite: http.SameSiteLaxMode,
					Expires:  time.Now().Add(duration),
					Path:     "/",
				},
			},
			ExpectedStatus: http.StatusOK,
			ExpectedJSON:   fmt.Sprintf(jsonResponseUserId, userIdString),
		},
		{
			Name:      "Wrong password",
			Method:    "POST",
			Route:     loginRoute,
			Env:       defaultEnv,
			InputJSON: fmt.Sprintf(jsonLoginUserWithUsername, username, pw),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows(rowsGetLogin).AddRow(userId, username, email, pwHashAlt))
			},
			ExpectedCookies: []*http.Cookie{},
			ExpectedStatus:  http.StatusUnauthorized,
			ExpectedError:   errors.Auth.WrongCredentialsError.ErrorMsg,
			ExpectedJSON:    jsonWrongCredentialsError,
		},
		{
			Name:      "Password missing",
			Method:    "POST",
			Route:     loginRoute,
			Env:       defaultEnv,
			InputJSON: fmt.Sprintf(jsonLoginUserWithUsernamePasswordMissing, username),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// should return before any database calls are made
			},
			ExpectedCookies: []*http.Cookie{},
			ExpectedStatus:  http.StatusUnauthorized,
			ExpectedError:   errors.Parameter.RequiredParameterMissingError.ErrorMsg,
			ExpectedJSON:    jsonPasswordRequiredError,
		},
		{
			Name:      "Both username and email missing",
			Method:    "POST",
			Route:     loginRoute,
			Env:       defaultEnv,
			InputJSON: fmt.Sprintf(jsonLoginUserUsernameAndEmailMissing, pw),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// should return before any database calls are made
			},
			ExpectedCookies: []*http.Cookie{},
			ExpectedStatus:  http.StatusUnauthorized,
			ExpectedError:   errors.Parameter.RequiredParameterMissingError.ErrorMsg,
			ExpectedJSON:    jsonUsernameOrEmailRequiredError,
		},
		{
			Name:      "Invalid payload",
			Method:    "POST",
			Route:     loginRoute,
			Env:       defaultEnv,
			InputJSON: fmt.Sprintf(jsonLoginUserWithUsernameInvalid, username, pw),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				// should return before any database calls are made
			},
			ExpectedCookies: []*http.Cookie{},
			ExpectedStatus:  http.StatusUnauthorized,
			ExpectedError:   errors.Payload.InvalidPayloadUserError.ErrorMsg,
			ExpectedJSON:    jsonInvalidUserPayloadError,
		},
		{
			Name:      "User not found (username)",
			Method:    "POST",
			Route:     loginRoute,
			Env:       defaultEnv,
			InputJSON: fmt.Sprintf(jsonLoginUserWithUsername, username, pw),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnRows(pgxmock.NewRows(rowsGetLogin))
			},
			ExpectedCookies: []*http.Cookie{},
			ExpectedStatus:  http.StatusUnauthorized,
			ExpectedError:   errors.Auth.WrongCredentialsError.ErrorMsg,
			ExpectedJSON:    jsonWrongCredentialsError,
		},
		{
			Name:      "User not found (email)",
			Method:    "POST",
			Route:     loginRoute,
			Env:       defaultEnv,
			InputJSON: fmt.Sprintf(jsonLoginUserWithEmail, email, pw),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetLogin))
			},
			ExpectedCookies: []*http.Cookie{},
			ExpectedStatus:  http.StatusUnauthorized,
			ExpectedError:   errors.Auth.WrongCredentialsError.ErrorMsg,
			ExpectedJSON:    jsonWrongCredentialsError,
		},
		{
			Name:      "JWT Secret missing",
			Method:    "POST",
			Route:     loginRoute,
			Env:       jwtSecretMissingEnv,
			InputJSON: fmt.Sprintf(jsonLoginUserWithEmail, email, pw),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetLogin))
			},
			ExpectedCookies: []*http.Cookie{},
			ExpectedStatus:  http.StatusUnauthorized,
			ExpectedError:   errors.Internal.JwtSecretMissingError.ErrorMsg,
		},
		{
			Name:      "Error during get query (email)",
			Method:    "POST",
			Route:     loginRoute,
			Env:       defaultEnv,
			InputJSON: fmt.Sprintf(jsonLoginUserWithEmail, email, pw),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnError(simulatedDbError)

			},
			ExpectedCookies: []*http.Cookie{},
			ExpectedStatus:  http.StatusInternalServerError,
			ExpectedJSON:    jsonSimulatedDbError,
		},
		{
			Name:      "Error during get query (username)",
			Method:    "POST",
			Route:     loginRoute,
			Env:       defaultEnv,
			InputJSON: fmt.Sprintf(jsonLoginUserWithUsername, username, pw),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByUsername).
					WithArgs(username).
					WillReturnError(simulatedDbError)
			},
			ExpectedCookies: []*http.Cookie{},
			ExpectedStatus:  http.StatusInternalServerError,
			ExpectedJSON:    jsonSimulatedDbError,
		},
		{
			Name:      "Error during set last logged in query",
			Method:    "POST",
			Route:     loginRoute,
			Env:       defaultEnv,
			InputJSON: fmt.Sprintf(jsonLoginUserWithEmail, email, pw),
			MockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(queryGetLoginByEmail).
					WithArgs(email).
					WillReturnRows(pgxmock.NewRows(rowsGetLogin).AddRow(userId, username, email, pwHash))

				mock.ExpectExec(queryUpdateLastLogin).WithArgs(userId).WillReturnError(simulatedDbError)
			},
			ExpectedCookies: []*http.Cookie{},
			ExpectedStatus:  http.StatusInternalServerError,
			ExpectedJSON:    jsonSimulatedDbError,
		},
	}
	testutil.RunTests(t, authHandler.HandleLogin, mock, tests)
}
