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
	"golang.org/x/crypto/bcrypt"
)

type createUserDto struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

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
	var userDto createUserDto

	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		log.Printf("Could not decode request body: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check the required parameters are present
	if len(userDto.Username) == 0 || len(userDto.Email) == 0 || len(userDto.Password) == 0 {
		log.Printf("One or more request parameters missing: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), 12)
	if err != nil {
		log.Printf("Could not hash password: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tx, err := dbConn.Begin(*dbCtx)
	if err != nil {
		log.Printf("Could not begin database transaction: %v\n", err)
	}
	defer tx.Rollback(*dbCtx)

	// converting the hash to a string here is not ideal, but sqlc interprets char(72) as a string so here we are
	createUserParams := db.CreateUserParams{Username: userDto.Username, Email: userDto.Email, Passwordhash: string(passwordHash)}
	userId, err := db.New(dbConn).WithTx(tx).CreateUser(*dbCtx, createUserParams)
	if err != nil {
		log.Printf("Could not create user: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	createUserDetailParams := db.CreateUserDetailsParams{ID: userId, Firstname: userDto.FirstName, Lastname: userDto.LastName}
	userId, err = db.New(dbConn).WithTx(tx).CreateUserDetails(*dbCtx, createUserDetailParams)
	if err != nil {
		log.Printf("Could not create user details: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// only commit the transaction once both the user and the user details have been created successfully
	err = tx.Commit(*dbCtx)
	if err != nil {
		log.Printf("Could not commit changes to database: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentType, applicationJson)
	err = json.NewEncoder(w).Encode(userId)
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
