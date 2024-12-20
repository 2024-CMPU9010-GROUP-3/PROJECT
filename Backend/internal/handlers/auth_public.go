//go:build public

package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	db "github.com/2024-CMPU9010-GROUP-3/magpie/internal/db/public"
	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/dtos"
	customErrors "github.com/2024-CMPU9010-GROUP-3/magpie/internal/errors"
	resp "github.com/2024-CMPU9010-GROUP-3/magpie/internal/responses"
	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/util/env"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

func (p *AuthHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	userIdPathParam := r.PathValue("id")
	var userId pgtype.UUID
	err := userId.Scan(userIdPathParam)

	// bad request if parameter is not valid uuid
	if err != nil {
		resp.SendError(customErrors.Parameter.InvalidUUIDError, w)
		return
	}

	userDetails, err := db.New(dbConn).GetUserDetails(*dbCtx, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			resp.SendError(customErrors.NotFound.UserNotFoundError, w)
			return
		} else {
			resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
			return
		}
	}

	resp.SendResponse(dtos.ResponseContentDto{
		HttpStatus: http.StatusOK,
		Content:    userDetails,
	}, w)
}

func (p *AuthHandler) HandlePost(w http.ResponseWriter, r *http.Request) {
	var userDto dtos.CreateUserDto

	e := userDto.Decode(r.Body)
	if e != nil {
		ce, ok := e.(customErrors.CustomError)
		if !ok {
			resp.SendError(customErrors.Internal.UnknownError.WithCause(e), w)
			return
		} else {
			resp.SendError(ce, w)
			return
		}
	}

	emailExists, err := db.New(dbConn).EmailExists(*dbCtx, db.EmailExistsParams{Email: userDto.Email, ID: pgtype.UUID{}})
	if err != nil {
		resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
		return
	}
	if emailExists {
		resp.SendError(customErrors.Payload.EmailAlreadyExistsError, w)
		return
	}

	usernameExists, err := db.New(dbConn).UsernameExists(*dbCtx, db.UsernameExistsParams{Username: userDto.Username, ID: pgtype.UUID{}})
	if err != nil {
		resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
		return
	}
	if usernameExists {
		resp.SendError(customErrors.Payload.UsernameAlreadyExistsError, w)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), 12)
	if err != nil {
		resp.SendError(customErrors.Internal.HashingError, w)
		return
	}

	tx, err := dbConn.Begin(*dbCtx)
	if err != nil {
		resp.SendError(customErrors.Database.TransactionStartError, w)
		return
	}
	defer func() {
		// potential error from rollback is not fatal, ignoring for now
		if tx != nil {
			_ = tx.Rollback(*dbCtx)
		}
	}()

	// converting the hash to a string here is not ideal, but sqlc interprets char(72) as a string so here we are
	createUserParams := db.CreateUserParams{Username: userDto.Username, Email: userDto.Email, Passwordhash: string(passwordHash)}
	userId, err := db.New(dbConn).WithTx(tx).CreateUser(*dbCtx, createUserParams)
	if err != nil {
		resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
		return
	}

	createUserDetailParams := db.CreateUserDetailsParams{ID: userId, Firstname: userDto.FirstName, Lastname: userDto.LastName, Profilepicture: userDto.ProfilePicture}
	userId, err = db.New(dbConn).WithTx(tx).CreateUserDetails(*dbCtx, createUserDetailParams)
	if err != nil {
		resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
		return
	}

	// only commit the transaction once both the user and the user details have been created successfully
	err = tx.Commit(*dbCtx)
	if err != nil {
		resp.SendError(customErrors.Database.TransactionCommitError, w)
		return
	}

	idDto := dtos.UserIdDto{UserId: userId}

	resp.SendResponse(dtos.ResponseContentDto{Content: idDto, HttpStatus: http.StatusCreated}, w)
}

func (p *AuthHandler) HandlePut(w http.ResponseWriter, r *http.Request) {
	var userId pgtype.UUID
	var userDto dtos.UpdateUserDto

	userIdPathParam := r.PathValue("id")

	err := userId.Scan(userIdPathParam)

	// bad request if parameter is not valid uuid
	if err != nil {
		resp.SendError(customErrors.Parameter.InvalidUUIDError, w)
		return
	}

	e := userDto.Decode(r.Body)
	if e != nil {
		ce, ok := e.(customErrors.CustomError)
		if !ok {
			resp.SendError(customErrors.Internal.UnknownError.WithCause(e), w)
			return
		} else {
			resp.SendError(ce, w)
			return
		}
	}

	emailExists, err := db.New(dbConn).EmailExists(*dbCtx, db.EmailExistsParams{Email: userDto.Email, ID: userId})
	if err != nil {
		resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
		return
	}
	if emailExists {
		resp.SendError(customErrors.Payload.EmailAlreadyExistsError, w)
		return
	}

	usernameExists, err := db.New(dbConn).UsernameExists(*dbCtx, db.UsernameExistsParams{Username: userDto.Username, ID: userId})
	if err != nil {
		resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
		return
	}
	if usernameExists {
		resp.SendError(customErrors.Payload.UsernameAlreadyExistsError, w)
		return
	}

	userLogin, err := db.New(dbConn).GetLoginById(*dbCtx, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			resp.SendError(customErrors.NotFound.UserNotFoundError, w)
			return
		} else {
			resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
			return
		}
	}

	var passwordHash []byte

	if userDto.Password != nil {
		passwordHash, err = bcrypt.GenerateFromPassword([]byte(*userDto.Password), 12)
		if err != nil {
			resp.SendError(customErrors.Internal.HashingError, w)
			return
		}
	} else {
		passwordHash = []byte(userLogin.Passwordhash)
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
		Profilepicture: userDto.ProfilePicture.String,
	}

	tx, err := dbConn.Begin(*dbCtx)
	if err != nil {
		resp.SendError(customErrors.Database.TransactionStartError, w)
		return
	}
	defer func() {
		// potential error from rollback is not fatal, ignoring for now
		if tx != nil {
			_ = tx.Rollback(*dbCtx)
		}
	}()

	err = db.New(dbConn).WithTx(tx).UpdateLogin(*dbCtx, updateLoginParams)
	if err != nil {
		resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
		return
	}

	err = db.New(dbConn).WithTx(tx).UpdateUserDetails(*dbCtx, updateUserDetailsParams)
	if err != nil {
		resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
		return
	}

	// only commit the transaction once both the user and the user details have been created successfully
	err = tx.Commit(*dbCtx)
	if err != nil {
		resp.SendError(customErrors.Database.TransactionCommitError, w)
		return
	}

	resp.SendResponse(dtos.ResponseContentDto{Content: dtos.UserIdDto{UserId: userId}, HttpStatus: http.StatusAccepted}, w)
}

func (p *AuthHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	var userId pgtype.UUID
	userIdPathParam := r.PathValue("id")
	err := userId.Scan(userIdPathParam)
	// bad request if parameter is not valid uuid
	if err != nil {
		resp.SendError(customErrors.Parameter.InvalidUUIDError, w)
		return
	}

	err = db.New(dbConn).DeleteUser(*dbCtx, userId)
	if err != nil {
		resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
		return
	}
	resp.SendResponse(dtos.ResponseContentDto{Content: dtos.UserIdDto{UserId: userId}, HttpStatus: http.StatusAccepted}, w)
}

// this method is very big and needs to be refactored in the future
func (p *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var loginDto dtos.UserLoginDto

	e := loginDto.Decode(r.Body)
	if e != nil {
		ce, ok := e.(customErrors.CustomError)
		if !ok {
			resp.SendError(customErrors.Internal.UnknownError.WithCause(e), w)
			return
		} else {
			resp.SendError(ce, w)
			return
		}
	}

	// get user login from db
	var userLogin db.Login
	var err error

	userLogin, err = db.New(dbConn).GetLoginByEmail(*dbCtx, loginDto.UsernameOrEmail)
	if err != nil {
		// try again with username
		userLogin, err = db.New(dbConn).GetLoginByUsername(*dbCtx, loginDto.UsernameOrEmail)
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			resp.SendError(customErrors.Auth.WrongCredentialsError, w)
			return
		} else {
			resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
			return
		}
	}

	// check password against password hash
	err = bcrypt.CompareHashAndPassword([]byte(userLogin.Passwordhash), []byte(loginDto.Password))
	if err != nil {
		resp.SendError(customErrors.Auth.WrongCredentialsError, w)
		return
	}

	// all env variables should be moved into a separate package and checked on startup in the future
	secret, set := env.Get(env.EnvJwtSecret)
	if !set {
		resp.SendError(customErrors.Internal.JwtSecretMissingError, w)
		return
	}
	expiry, _ := env.Get(env.EnvJwtExpiry)

	parsedExpiry, err := time.ParseDuration(expiry)
	if err != nil {
		resp.SendError(customErrors.Internal.UnknownError.WithCause(fmt.Errorf("Could not parse JWT expiry duration")), w)
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
		resp.SendError(customErrors.Internal.UnknownError.WithCause(fmt.Errorf("Could not sign JWT")), w)
		return
	}

	// set last logged in in database
	err = db.New(dbConn).UpdateLastLogin(*dbCtx, userLogin.ID)
	if err != nil {
		resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
		return
	}

	tokenDto := dtos.UserLoginResponseDto{
		UserId: userLogin.ID,
		Token:  tokenString,
	}

	resp.SendResponse(dtos.ResponseContentDto{Content: tokenDto, HttpStatus: http.StatusOK}, w)
}
