// This will need refactoring in the future, but it is sufficient for the amount of routes we have at the moment
package routes

import (
	"net/http"

	"github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/handlers"
)

func Router() *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("/v1/", http.StripPrefix("/v1", v1()))
	return router
}

func v1() *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("/private/", http.StripPrefix("/private", private()))
	router.Handle("/public/", http.StripPrefix("/public", public()))
	return router
}

func private() *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("/points/", http.StripPrefix("/points", pointsPrivate()))
	return router

}

func public() *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("/points/", http.StripPrefix("/points", pointsPublic()))
	router.Handle("/auth/", http.StripPrefix("/auth", auth()))
	return router
}

func pointsPrivate() *http.ServeMux {
	router := http.NewServeMux()
	pointsHandler := &handlers.PointsHandler{}
	router.HandleFunc("POST /", pointsHandler.HandlePost)
	router.HandleFunc("PUT /{id}", pointsHandler.HandlePut)
	router.HandleFunc("DELETE /{id}", pointsHandler.HandleDelete)

	return router
}

func pointsPublic() *http.ServeMux {
	router := http.NewServeMux()
	pointsHandler := &handlers.PointsHandler{}
	router.HandleFunc("GET /byRadius", pointsHandler.HandleGetByRadius)
	router.HandleFunc("GET /pointDetails", pointsHandler.HandleGetPointDetails)

	return router
}

func auth() *http.ServeMux {
	router := http.NewServeMux()
	authHandler := &handlers.AuthHandler{}

	router.HandleFunc("GET /User/{id}", authHandler.HandleGet)
	router.HandleFunc("POST /User", authHandler.HandlePost)
	router.HandleFunc("PUT /User/{id}", authHandler.HandlePut)
	router.HandleFunc("DELETE /User/{id}", authHandler.HandleDelete)
	router.HandleFunc("POST /User/login", authHandler.HandleLogin)

	return router
}
