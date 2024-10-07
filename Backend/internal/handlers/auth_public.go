//go:build public

package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	db "github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/db/public"
	"github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (p *AuthHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	userIdPathParam := r.PathValue("id")
	log.Printf("User id path param: %v\n", userIdPathParam)
	var userId pgtype.UUID
	err := userId.Scan(userIdPathParam)

	// bad request if parameter is not valid uuid
	if err != nil {
		log.Printf("Invalid path parameter: %v\n", userIdPathParam)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userDetails, err := db.New(dbConn).GetUserDetails(*dbCtx, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("User details not found in database: %+v\n", err)
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			log.Printf("Could not get user details from database: %+v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set(contentType, applicationJson)
	err = json.NewEncoder(w).Encode(userDetails)
}

func (p *AuthHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	data := util.Placeholder("POST User")
	w.Header().Set(contentType, applicationJson)
	err := json.NewEncoder(w).Encode(data)
	util.CheckResponseError(err, w)
}

func (p *AuthHandler) HandlePut(w http.ResponseWriter, r *http.Request) {
	data := util.Placeholder("PUT User")
	w.Header().Set(contentType, applicationJson)
	err := json.NewEncoder(w).Encode(data)
	util.CheckResponseError(err, w)
}

func (p *AuthHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	data := util.Placeholder("DELETE User")
	w.Header().Set(contentType, applicationJson)
	err := json.NewEncoder(w).Encode(data)
	util.CheckResponseError(err, w)
}

func (p *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	data := util.Placeholder("POST User/Login")
	w.Header().Set(contentType, applicationJson)
	err := json.NewEncoder(w).Encode(data)
	util.CheckResponseError(err, w)
}
