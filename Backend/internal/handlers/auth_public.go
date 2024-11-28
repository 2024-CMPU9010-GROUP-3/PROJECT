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
	util "github.com/2024-CMPU9010-GROUP-3/magpie/internal/util"
	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/util/env"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

func (p *AuthHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	userId, err := util.GetUserIdFromRequest(r)
	if err != nil {
		resp.SendError(err.(customErrors.CustomError), w)
		return
	}

	userDetails, err := PublicDb().GetUserDetails(*dbCtx, userId)
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

	if err := p.checkForConflicts(pgtype.UUID{}, userDto.Email, userDto.Username); err != nil {
		resp.SendError(err.(customErrors.CustomError), w)
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
	userId, err := PublicDb().WithTx(tx).CreateUser(*dbCtx, createUserParams)
	if err != nil {
		resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
		return
	}

	createUserDetailParams := db.CreateUserDetailsParams{ID: userId, Firstname: userDto.FirstName, Lastname: userDto.LastName, Profilepicture: userDto.ProfilePicture}
	userId, err = PublicDb().WithTx(tx).CreateUserDetails(*dbCtx, createUserDetailParams)
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
	var userDto dtos.UpdateUserDto
	userId, err := util.GetUserIdFromRequest(r)
	if err != nil {
		resp.SendError(err.(customErrors.CustomError), w)
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

	if err := p.checkForConflicts(userId, userDto.Email, userDto.Username); err != nil {
		resp.SendError(err.(customErrors.CustomError), w)
		return
	}

	userLogin, err := PublicDb().GetLoginById(*dbCtx, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			resp.SendError(customErrors.NotFound.UserNotFoundError, w)
			return
		} else {
			resp.SendError(customErrors.Database.UnknownDatabaseError.WithCause(err), w)
			return
		}
	}

	passwordHash, err := p.getPasswordHash(userDto.Password, userLogin.Passwordhash)
	if err != nil {
		resp.SendError(err.(customErrors.CustomError), w)
		return
	}

	err = p.updateUser(userId, userDto, passwordHash)
	if err != nil {
		resp.SendError(err.(customErrors.CustomError), w)
		return
	}

	resp.SendResponse(dtos.ResponseContentDto{Content: dtos.UserIdDto{UserId: userId}, HttpStatus: http.StatusAccepted}, w)
}

func (p *AuthHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	userId, err := util.GetUserIdFromRequest(r)
	if err != nil {
		resp.SendError(err.(customErrors.CustomError), w)
		return
	}

	err = PublicDb().DeleteUser(*dbCtx, userId)
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

	userLogin, err = PublicDb().GetLoginByEmail(*dbCtx, loginDto.UsernameOrEmail)
	if err != nil {
		// try again with username
		userLogin, err = PublicDb().GetLoginByUsername(*dbCtx, loginDto.UsernameOrEmail)
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

	tokenString, err := createJWTToken(userLogin.ID)
	if err != nil {
		resp.SendError(err.(customErrors.CustomError), w)
		return
	}

	// set last logged in in database
	err = PublicDb().UpdateLastLogin(*dbCtx, userLogin.ID)
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

// Helper methods

func (p *AuthHandler) checkForConflicts(userId pgtype.UUID, email string, username string) error {
	if exists, err := PublicDb().EmailExists(*dbCtx, db.EmailExistsParams{Email: email, ID: userId}); err != nil {
		return customErrors.Database.UnknownDatabaseError.WithCause(err)
	} else if exists {
		return customErrors.Payload.EmailAlreadyExistsError
	}

	if exists, err := PublicDb().UsernameExists(*dbCtx, db.UsernameExistsParams{Username: username, ID: userId}); err != nil {
		return customErrors.Database.UnknownDatabaseError.WithCause(err)
	} else if exists {
		return customErrors.Payload.UsernameAlreadyExistsError
	}
	return nil
}

func (p *AuthHandler) getPasswordHash(newPassword *string, existingHash string) ([]byte, error) {
	if newPassword == nil {
		return []byte(existingHash), nil
	}
	return bcrypt.GenerateFromPassword([]byte(*newPassword), 12)
}

func (p *AuthHandler) updateUser(userId pgtype.UUID, userDto dtos.UpdateUserDto, passwordHash []byte) error {
	return withTransaction(func(tx pgx.Tx) error {
		// Update login data
		if err := db.New(dbConn).WithTx(tx).UpdateLogin(*dbCtx, db.UpdateLoginParams{
			ID:           userId,
			Username:     userDto.Username,
			Email:        userDto.Email,
			Passwordhash: string(passwordHash),
		}); err != nil {
			return customErrors.Database.UnknownDatabaseError.WithCause(err)
		}

		// Update user details
		if err := db.New(dbConn).WithTx(tx).UpdateUserDetails(*dbCtx, db.UpdateUserDetailsParams{
			ID:             userId,
			Firstname:      userDto.FirstName,
			Lastname:       userDto.LastName,
			Profilepicture: userDto.ProfilePicture.String,
		}); err != nil {
			return customErrors.Database.UnknownDatabaseError.WithCause(err)
		}

		return nil
	})
}

func createJWTToken(userId pgtype.UUID) (string, error) {
	secret, set := env.Get(env.EnvJwtSecret)
	if !set {
		return "", customErrors.Internal.JwtSecretMissingError
	}
	expiry, _ := env.Get(env.EnvJwtExpiry)
	parsedExpiry, err := time.ParseDuration(expiry)
	if err != nil {
		return "", customErrors.Internal.UnknownError.WithCause(fmt.Errorf("Could not parse JWT expiry duration: %w", err))
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(parsedExpiry).Unix(),
		"iat": time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		err = customErrors.Internal.UnknownError.WithCause(fmt.Errorf("Could not sign JWT"))
	}
	return tokenString, err
}
