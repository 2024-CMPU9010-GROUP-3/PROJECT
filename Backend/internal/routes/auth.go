package routes

import "net/http"

type authRoute struct{}

func (a *authRoute) Router() *http.ServeMux {
	router := http.NewServeMux()

	return router
}
