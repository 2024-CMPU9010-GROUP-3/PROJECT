package routes

import "net/http"

type privateRoute struct{}

func (p *privateRoute) Router() *http.ServeMux {
	pointsRoute := &pointsRoute{}
	router := http.NewServeMux()

	router.Handle("/points/", http.StripPrefix("/points", pointsRoute.PrivateRouter()))

	return router
}
