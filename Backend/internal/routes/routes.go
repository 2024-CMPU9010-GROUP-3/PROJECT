package routes

import (
	"net/http"
)

func Router() *http.ServeMux {
	router := http.NewServeMux()
	v1Route := &v1Route{}

	router.Handle("/v1/", http.StripPrefix("/v1", v1Route.Router()))
	return router
}
