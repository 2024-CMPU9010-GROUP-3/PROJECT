package routes

import "net/http"

func LoadRouter(router *http.ServeMux) {
	v1Route := &v1Route{}

	router.Handle("/v1/", http.StripPrefix("/v1", v1Route.Router()))
}
