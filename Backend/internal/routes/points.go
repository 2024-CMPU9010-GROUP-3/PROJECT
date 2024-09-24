package routes

import "net/http"

type pointsRoute struct{}

func (p *pointsRoute) PublicRouter() *http.ServeMux {
	router := http.NewServeMux()

	return router
}

func (p *pointsRoute) PrivateRouter() *http.ServeMux {
	router := http.NewServeMux()

	return router
}
