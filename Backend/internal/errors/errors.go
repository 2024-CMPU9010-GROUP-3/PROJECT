package errors

type CustomError struct {
	HttpStatus int
	ErrorCode  int    `json:"errorCode"`
	ErrorMsg   string `json:"errorMsg"`
}
