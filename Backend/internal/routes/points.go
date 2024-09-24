package routes

import (
	"net/http"

	"github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/handlers"
)

type pointsRoute struct{}

func (p *pointsRoute) PublicRouter() *http.ServeMux {
	router := http.NewServeMux()
	pointsHandler := &handlers.PointsHandler{}
	router.HandleFunc("GET /byRadius", pointsHandler.HandleGetByRadius)
	router.HandleFunc("GET /pointDetails", pointsHandler.HandleGetPointDetails)

	return router
}

func (p *pointsRoute) PrivateRouter() *http.ServeMux {
	router := http.NewServeMux()
	pointsHandler := &handlers.PointsHandler{}
	router.HandleFunc("POST /", pointsHandler.HandlePost)
	router.HandleFunc("PUT /", pointsHandler.HandlePut)
	router.HandleFunc("DELETE /", pointsHandler.HandleDelete)

	return router
}
