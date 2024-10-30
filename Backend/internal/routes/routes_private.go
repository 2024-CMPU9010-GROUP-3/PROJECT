//go:build private

// This will need refactoring in the future, but it is sufficient for the amount of routes we have at the moment

package routes

import (
	"net/http"

	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/handlers"
)

func init() {
	AddRoute(route{"/private/", private()})
}

func private() *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("/points/", http.StripPrefix("/points", pointsPrivate()))
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
