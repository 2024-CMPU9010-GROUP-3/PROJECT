package errors

import "fmt"

type CustomError struct {
	HttpStatus int    `json:"-"`
	ErrorCode  int    `json:"errorCode"`
	ErrorMsg   string `json:"errorMsg"`
	Cause      string `json:"cause,omitempty"`
}

func (c CustomError) Error() string {
	if c.Cause != ""{
		return fmt.Sprintf("%d: %s (%s)", c.ErrorCode, c.ErrorMsg, c.Cause)
	} else {
		return fmt.Sprintf("%d: %s", c.ErrorCode, c.ErrorMsg)
	}
}

func (c CustomError) WithCause(err error) CustomError {
	c.Cause = err.Error()
	return c
}
