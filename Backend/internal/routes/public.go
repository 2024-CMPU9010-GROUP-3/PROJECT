package routes

import "net/http"

type publicRoute struct{}

func (p *publicRoute) Router() *http.ServeMux {
	pointsRoute := &pointsRoute{}
	authRoute := &authRoute{}
	router := http.NewServeMux()

	router.Handle("/points/", http.StripPrefix("/points", pointsRoute.PublicRouter()))
	router.Handle("/auth/", http.StripPrefix("/auth", authRoute.Router()))

	return router
}
