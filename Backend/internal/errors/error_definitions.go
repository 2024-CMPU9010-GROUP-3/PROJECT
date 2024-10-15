package errors

import "net/http"

const (
	codeUnknownError          = 1001
	codeJwtSecretMissingError = 1002
	codeHashingError          = 1003
	codeJsonEncodingError     = 1004

	codeDatabaseConnectionError        = 1101
	codeDatabaseTransactionStartError  = 1102
	codeDatabaseTransactionCommitError = 1103
	codeUnknownDatabaseError           = 1104

	codeInvalidParameterUUIDError      = 1201
	codeInvalidParameterFloatError     = 1202
	codeInvalidParameterLongitudeError = 1203
	codeInvalidParameterLatitudeError  = 1204
	codeInvalidPayloadPointError       = 1205
	codeInvalidPayloadUserError        = 1206
	codeUsernameAlreadyExistsError     = 1207
	codeEmailAlreadyExistsError        = 1208
	codeRequiredParameterMissingError  = 1209

	codeUserNotFoundError  = 1301
	codePointNotFoundError = 1302

	codeUnauthorizedError     = 1401
	codeWrongCredentialsError = 1402
	codeInvalidTokenError     = 1403
)

var unknownError = CustomError{
	HttpStatus: http.StatusInternalServerError,
	ErrorCode:  codeUnknownError,
	ErrorMsg:   "Unknown internal server error",
}

var jwtSecretMissingError = CustomError{
	HttpStatus: http.StatusInternalServerError,
	ErrorCode:  codeJwtSecretMissingError,
	ErrorMsg:   "Could not generate JWT, secret not set",
}

var hashingError = CustomError{
	HttpStatus: http.StatusInternalServerError,
	ErrorCode:  codeHashingError,
	ErrorMsg:   "Could not hash password",
}

var jsonEncodingError = CustomError{
	HttpStatus: http.StatusInternalServerError,
	ErrorCode:  codeJsonEncodingError,
	ErrorMsg:   "Could not encode json",
}

var databaseConnectionError = CustomError{
	HttpStatus: http.StatusInternalServerError,
	ErrorCode:  codeDatabaseConnectionError,
	ErrorMsg:   "Could not connect to database",
}

var databaseTransactionStartError = CustomError{
	HttpStatus: http.StatusInternalServerError,
	ErrorCode:  codeDatabaseTransactionStartError,
	ErrorMsg:   "Could not start database transaction",
}

var databaseTransactionCommitError = CustomError{
	HttpStatus: http.StatusInternalServerError,
	ErrorCode:  codeDatabaseTransactionCommitError,
	ErrorMsg:   "Could not commit database transaction",
}

var databaseUnknownError = CustomError{
	HttpStatus: http.StatusInternalServerError,
	ErrorCode:  codeUnknownDatabaseError,
	ErrorMsg:   "Unknown database error",
}

var invalidParameterUUIDError = CustomError{
	HttpStatus: http.StatusBadRequest,
	ErrorCode:  codeInvalidParameterUUIDError,
	ErrorMsg:   "Parameter invalid, expected type UUIDv4",
}

var invalidParameterFloatError = CustomError{
	HttpStatus: http.StatusBadRequest,
	ErrorCode:  codeInvalidParameterFloatError,
	ErrorMsg:   "Parameter invalid, expected type Float",
}

var invalidParameterLongitudeError = CustomError{
	HttpStatus: http.StatusBadRequest,
	ErrorCode:  codeInvalidParameterLongitudeError,
	ErrorMsg:   "Parameter invalid, expected Longitude",
}

var invalidParameterLatitudeError = CustomError{
	HttpStatus: http.StatusBadRequest,
	ErrorCode:  codeInvalidParameterLatitudeError,
	ErrorMsg:   "Parameter invalid, expected Latitude",
}

var invalidPayloadPointError = CustomError{
	HttpStatus: http.StatusBadRequest,
	ErrorCode:  codeInvalidPayloadPointError,
	ErrorMsg:   "Payload (Point) not valid",
}

var invalidPayloadUserError = CustomError{
	HttpStatus: http.StatusBadRequest,
	ErrorCode:  codeInvalidPayloadUserError,
	ErrorMsg:   "Payload (User) not valid",
}

var usernameAlreadyExistsError = CustomError{
	HttpStatus: http.StatusBadRequest,
	ErrorCode:  codeUsernameAlreadyExistsError,
	ErrorMsg:   "Username already exists",
}

var emailAlreadyExistsError = CustomError{
	HttpStatus: http.StatusBadRequest,
	ErrorCode:  codeEmailAlreadyExistsError,
	ErrorMsg:   "Email already exists",
}

var requiredParameterMissingError = CustomError{
	HttpStatus: http.StatusBadRequest,
	ErrorCode:  codeRequiredParameterMissingError,
	ErrorMsg:   "One or more required parameters are missing",
}

var userNotFoundError = CustomError{
	HttpStatus: http.StatusNotFound,
	ErrorCode:  codeUserNotFoundError,
	ErrorMsg:   "User not found",
}

var pointNotFoundError = CustomError{
	HttpStatus: http.StatusNotFound,
	ErrorCode:  codePointNotFoundError,
	ErrorMsg:   "Point not found",
}

var unauthorizedError = CustomError{
	HttpStatus: http.StatusUnauthorized,
	ErrorCode:  codeUnauthorizedError,
	ErrorMsg:   "Unauthorized",
}

var wrongCredentialsError = CustomError{
	HttpStatus: http.StatusUnauthorized,
	ErrorCode:  codeWrongCredentialsError,
	ErrorMsg:   "Wrong username or password",
}

var invalidTokenError = CustomError{
	HttpStatus: http.StatusUnauthorized,
	ErrorCode:  codeInvalidTokenError,
	ErrorMsg:   "Bearer token not valid",
}

var Internal = struct {
	UnknownError          CustomError
	JwtSecretMissingError CustomError
	HashingError          CustomError
	JsonEncodingError     CustomError
}{
	UnknownError:          unknownError,
	JwtSecretMissingError: jwtSecretMissingError,
	HashingError:          hashingError,
	JsonEncodingError:     jsonEncodingError,
}

var Database = struct {
	ConnectionError        CustomError
	TransactionStartError  CustomError
	TransactionCommitError CustomError
	UnknownDatabaseError   CustomError
}{
	ConnectionError:        databaseConnectionError,
	TransactionStartError:  databaseTransactionStartError,
	TransactionCommitError: databaseTransactionCommitError,
	UnknownDatabaseError:   databaseUnknownError,
}

var Parameter = struct {
	InvalidUUIDError      CustomError
	InvalidFloatError     CustomError
	InvalidLongitudeError CustomError
	InvalidLatitudeError  CustomError
}{
	InvalidUUIDError:      invalidParameterUUIDError,
	InvalidFloatError:     invalidParameterFloatError,
	InvalidLongitudeError: invalidParameterLongitudeError,
	InvalidLatitudeError:  invalidParameterLatitudeError,
}

var Payload = struct {
	InvalidPayloadPointError      CustomError
	InvalidPayloadUserError       CustomError
	UsernameAlreadyExistsError    CustomError
	EmailAlreadyExistsError       CustomError
	RequiredParameterMissingError CustomError
}{
	InvalidPayloadPointError:      invalidPayloadPointError,
	InvalidPayloadUserError:       invalidPayloadUserError,
	UsernameAlreadyExistsError:    usernameAlreadyExistsError,
	EmailAlreadyExistsError:       emailAlreadyExistsError,
	RequiredParameterMissingError: requiredParameterMissingError,
}

var NotFound = struct {
	UserNotFoundError  CustomError
	PointNotFoundError CustomError
}{
	UserNotFoundError:  userNotFoundError,
	PointNotFoundError: pointNotFoundError,
}

var Auth = struct {
	UnauthorizedError     CustomError
	WrongCredentialsError CustomError
	InvalidTokenError     CustomError
}{
	UnauthorizedError:     unauthorizedError,
	WrongCredentialsError: wrongCredentialsError,
	InvalidTokenError:     invalidTokenError,
}
