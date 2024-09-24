package routes

import (
	"net/http"

	"github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/handlers"
)

type authRoute struct{}

func (a *authRoute) Router() *http.ServeMux {
	router := http.NewServeMux()
	authHandler := &handlers.AuthHandler{}

	router.HandleFunc("GET /User", authHandler.HandleGet)
	router.HandleFunc("POST /User", authHandler.HandlePost)
	router.HandleFunc("PUT /User", authHandler.HandlePut)
	router.HandleFunc("DELETE /User", authHandler.HandleDelete)
	router.HandleFunc("POST /User/login", authHandler.HandleLogin)

	return router
}
