package responses

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/dtos"
	errs "github.com/2024-CMPU9010-GROUP-3/magpie/internal/errors"
)

func SendResponse(response dtos.ResponseContentDto, w http.ResponseWriter) {
	resp := dtos.ResponseDto{
		Error:    nil,
		Response: &response,
	}
	w.WriteHeader(response.HttpStatus)
	send(resp, w)
}

func SendError(err errs.CustomError, w http.ResponseWriter) {
	resp := dtos.ResponseDto{
		Error:    &err,
		Response: nil,
	}
	log.Printf("Error: %+v\n", err)
	w.WriteHeader(err.HttpStatus)
	send(resp, w)
}

func send(resp dtos.ResponseDto, w http.ResponseWriter) {
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		// In this case we cannot send any data to the frontend
		log.Printf("Could not send response: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
