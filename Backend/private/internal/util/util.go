package util

import (
	"fmt"
	"net/http"
)

type placeholder struct {
	IsPlaceholder bool
	Endpoint      string
}

func Placeholder(endpoint string) *placeholder {
	return &placeholder{true, endpoint}
}

func CheckResponseError(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, fmt.Sprintf("response error, %v", err), http.StatusInternalServerError)
	}
}