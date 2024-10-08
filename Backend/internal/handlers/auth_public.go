//go:build public

package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	db "github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/db/public"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

const secretEnv = "MAGPIE_JWT_SECRET"
const expiryEnv = "MAGPIE_JWT_EXPIRY"

// these dtos need to be refactored into their own package in the future
type createUserDto struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	FirstName      string `json:"firstname"`
	LastName       string `json:"lastname"`
	ProfilePicture string `json:"profilepicture"`
}

type userLoginDto struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type bearerTokenDto struct {
	BearerToken string `json:"bearertoken"`
}

type userIdDto struct {
	UserId pgtype.UUID `json:"userid"`
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

	err = json.NewEncoder(w).Encode(userDetails)
	if err != nil {
		log.Printf("Could not send user details as response: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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

	idDto := userIdDto{UserId: userId}

	err = json.NewEncoder(w).Encode(idDto)
	if err != nil {
		log.Printf("Could not send user id as response: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p *AuthHandler) HandlePut(w http.ResponseWriter, r *http.Request) {
	var userId pgtype.UUID
	var userDto createUserDto

	userIdPathParam := r.PathValue("id")

	err := userId.Scan(userIdPathParam)

	// bad request if parameter is not valid uuid
	if err != nil {
		log.Printf("Invalid path parameter: %v\n", userIdPathParam)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		log.Printf("Could not decode request body: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var passwordHash []byte

	if len(userDto.Password) != 0 {
		passwordHash, err = bcrypt.GenerateFromPassword([]byte(userDto.Password), 12)
		if err != nil {
			log.Printf("Could not hash password: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	updateLoginParams := db.UpdateLoginParams{
		ID:           userId,
		Username:     userDto.Username,
		Email:        userDto.Email,
		Passwordhash: string(passwordHash),
	}

	updateUserDetailsParams := db.UpdateUserDetailsParams{
		ID:             userId,
		Firstname:      userDto.FirstName,
		Lastname:       userDto.LastName,
		Profilepicture: userDto.ProfilePicture,
	}

	log.Printf("update login params: %+v\n", updateLoginParams)
	log.Printf("update user details params: %+v\n", updateUserDetailsParams)
	tx, err := dbConn.Begin(*dbCtx)
	if err != nil {
		log.Printf("Could not begin database transaction: %v\n", err)
	}
	defer tx.Rollback(*dbCtx)

	err = db.New(dbConn).WithTx(tx).UpdateLogin(*dbCtx, updateLoginParams)
	if err != nil {
		log.Printf("Could not update login: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = db.New(dbConn).WithTx(tx).UpdateUserDetails(*dbCtx, updateUserDetailsParams)
	if err != nil {
		log.Printf("Could not update user details: %+v\n", err)
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
}

func (p *AuthHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	var userId pgtype.UUID
	userIdPathParam := r.PathValue("id")
	err := userId.Scan(userIdPathParam)
	// bad request if parameter is not valid uuid
	if err != nil {
		log.Printf("Invalid path parameter: %v\n", userIdPathParam)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = db.New(dbConn).DeleteUser(*dbCtx, userId)
	if err != nil {
		log.Printf("Could not delete user from database: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// this method is very big and needs to be refactored in the future
func (p *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var loginDto userLoginDto
	err := json.NewDecoder(r.Body).Decode(&loginDto)
	if err != nil {
		log.Printf("Could not decode request body: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(loginDto.Password) == 0 {
		log.Printf("Password is a required parameter\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(loginDto.Username) == 0 {
		log.Printf("Username is a required parameter\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get user login from db
	var userLogin db.Login
	userLogin, err = db.New(dbConn).GetLoginByUsername(*dbCtx, loginDto.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("User not found\n")
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			log.Printf("Could not get login details from database: %+v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// check password against password hash
	err = bcrypt.CompareHashAndPassword([]byte(userLogin.Passwordhash), []byte(loginDto.Password))
	if err != nil {
		log.Printf("Password incorrect: %+v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// all env variables should be moved into a separate package and checked on startup in the future
	secret := os.Getenv(secretEnv)
	if len(secret) == 0 {
		log.Printf("Could not generate JWT, environment variable MAGPIE_JWT_SECRET not set")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	expiry := os.Getenv(expiryEnv)
	if len(secret) == 0 {
		log.Printf("Warning: MAGPIE_JWT_EXPIRY not set, defaulting to 7d expiry")
		expiry = "24h"
	}

	parsedExpiry, err := time.ParseDuration(expiry)
	if err != nil {
		log.Printf("Could not parse JWT expiry to duration: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// generate bearer token with user id in payload
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userLogin.ID,
		"exp": time.Now().Add(parsedExpiry).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Printf("Could not generate signed JWT: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// set last logged in in database
	err = db.New(dbConn).UpdateLastLogin(*dbCtx, userLogin.ID)
	if err != nil {
		log.Printf("Could not update login date in database: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tokenDto := bearerTokenDto{
		tokenString,
	}

	// send bearer token as response
	err = json.NewEncoder(w).Encode(tokenDto)
	if err != nil {
		log.Printf("Could not send bearer token as response: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
