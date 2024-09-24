package routes

import (
	"net/http"
)

type v1Route struct{}

func (v *v1Route) Router() *http.ServeMux {
	publicHandler := &publicRoute{}
	privateHandler := &privateRoute{}
	router := http.NewServeMux()

	router.Handle("/private/", http.StripPrefix("/private", privateHandler.Router()))
	router.Handle("/public/", http.StripPrefix("/public", publicHandler.Router()))

	return router
}
