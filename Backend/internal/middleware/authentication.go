package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func accessAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		auth_cookie, err := request.Cookie("magpie_auth")
		if err != nil {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(auth_cookie.Value, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			jwtSecret := []byte(os.Getenv("MAGPIE_JWT_SECRET"))
			if len(jwtSecret) == 0 {
				return nil, fmt.Errorf("Cannot decrypt JWT, MAGPIE_JWT_SECRET not set")
			}

			return jwtSecret, nil
		})
		if err != nil {
			log.Printf("Auth token could not be parsed: %+v", err)
			return
		}

		subject, err := token.Claims.GetSubject()
		if err != nil {
			log.Printf("Could not extract user id from auth token: %+v\n", err)
			return
		}

		ctx := context.WithValue(request.Context(), "token_user_id", subject)

		next.ServeHTTP(writer, request.WithContext(ctx))

	})
}

func accessOwnerOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		userIdPathParam := request.PathValue("id")
		var userId pgtype.UUID
		err := userId.Scan(userIdPathParam)

		// bad request if parameter is not valid uuid
		if err != nil {
			log.Printf("Invalid value for user id")
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		tokenUserId := request.Context().Value("token_user_id")
		if tokenUserId == nil {
			log.Printf("No user_id from token present in context")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if tokenUserId != userIdPathParam {
			log.Printf("User not authorised for this request")
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(writer, request)
	})
}

func accessPublic(next http.Handler) http.Handler {
	// this is a nothing-wrapper, but explicit is better than implicit
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		next.ServeHTTP(writer, request)
	})
}
