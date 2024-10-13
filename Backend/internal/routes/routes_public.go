//go:build public

// This will need refactoring in the future, but it is sufficient for the amount of routes we have at the moment

package routes

import (
	"net/http"

	"github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/handlers"
)

// @title Public API
// @version 1.0
// @description This is the public API for the project.
// @BasePath /public

func init() {
	AddRoute(route{"/public/", public()})
}

func public() *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("/points/", http.StripPrefix("/points", pointsPublic()))
	router.Handle("/auth/", http.StripPrefix("/auth", auth()))

	return router
}

func pointsPublic() *http.ServeMux {
	router := http.NewServeMux()
	pointsHandler := &handlers.PointsHandler{}

	// Write code here for swagger annotations for the below path
	// swagger:operation GET /product Product getList
	// Get Point by Radius
	//
	// ---
	// responses:
	//
	//  401: CommonError
	//  200: CommonSuccess
	router.HandleFunc("GET /byRadius", pointsHandler.HandleGetByRadius)

	// @Summary Get point details
	// @Description Get details of a specific point by ID
	// @Tags points
	// @Accept  json
	// @Produce  json
	// @Param   id  path  string  true  "Point ID"
	// @Success 200 {object} handlers.PointDetailsResponse
	// @Router /points/{id} [get]
	router.HandleFunc("GET /{id}", pointsHandler.HandleGetPointDetails)

	return router
}

func auth() *http.ServeMux {
	router := http.NewServeMux()
	authHandler := &handlers.AuthHandler{}

	// @Summary Get user details
	// @Description Get details of a specific user by ID
	// @Tags auth
	// @Accept  json
	// @Produce  json
	// @Param   id  path  string  true  "User ID"
	// @Success 200 {object} handlers.UserResponse
	// @Router /auth/User/{id} [get]
	router.HandleFunc("GET /User/{id}", authHandler.HandleGet)

	// @Summary Create a new user
	// @Description Create a new user
	// @Tags auth
	// @Accept  json
	// @Produce  json
	// @Param   user  body  handlers.UserRequest  true  "User data"
	// @Success 201 {object} handlers.UserResponse
	// @Router /auth/User [post]
	router.HandleFunc("POST /User/", authHandler.HandlePost)

	// @Summary Update user details
	// @Description Update details of a specific user by ID
	// @Tags auth
	// @Accept  json
	// @Produce  json
	// @Param   id  path  string  true  "User ID"
	// @Param   user  body  handlers.UserRequest  true  "User data"
	// @Success 200 {object} handlers.UserResponse
	// @Router /auth/User/{id} [put]
	router.HandleFunc("PUT /User/{id}", authHandler.HandlePut)

	// @Summary Delete a user
	// @Description Delete a specific user by ID
	// @Tags auth
	// @Accept  json
	// @Produce  json
	// @Param   id  path  string  true  "User ID"
	// @Success 204
	// @Router /auth/User/{id} [delete]
	router.HandleFunc("DELETE /User/{id}", authHandler.HandleDelete)

	// @Summary User login
	// @Description Authenticate a user
	// @Tags auth
	// @Accept  json
	// @Produce  json
	// @Param   credentials  body  handlers.LoginRequest  true  "Login credentials"
	// @Success 200 {object} handlers.LoginResponse
	// @Router /auth/User/login [post]
	router.HandleFunc("POST /User/login", authHandler.HandleLogin)

	return router
}
