package errors

type CustomError struct {
	HttpStatus int
	ErrorCode  int    `json:"errorCode"`
	ErrorMsg   string `json:"errorMsg"`
	Cause      string `json:"cause"`
}

func (c CustomError) WithCause(err error) CustomError {
	c.Cause = err.Error()
	return c
}
