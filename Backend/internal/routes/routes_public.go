//go:build public

// This will need refactoring in the future, but it is sufficient for the amount of routes we have at the moment

package routes

import (
	"net/http"

	"github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/handlers"
	"github.com/go-openapi/runtime/middleware"
)

func init() {
	AddRoute(route{"/public/", public()})
}

func public() *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("/points/", http.StripPrefix("/points", pointsPublic()))
	router.Handle("/auth/", http.StripPrefix("/auth", auth()))

	opts := middleware.SwaggerUIOpts{SpecURL: "./swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	router.Handle("/docs", sh)
	router.HandleFunc("/swagger.yaml", swaggerSpec)

	return router
}

func pointsPublic() *http.ServeMux {
	router := http.NewServeMux()
	pointsHandler := &handlers.PointsHandler{}
	router.HandleFunc("GET /byRadius", pointsHandler.HandleGetByRadius)
	router.HandleFunc("GET /{id}", pointsHandler.HandleGetPointDetails)

	return router
}

func auth() *http.ServeMux {
	router := http.NewServeMux()
	authHandler := &handlers.AuthHandler{}

	router.HandleFunc("GET /User/{id}", authHandler.HandleGet)
	router.HandleFunc("POST /User/", authHandler.HandlePost)
	router.HandleFunc("PUT /User/{id}", authHandler.HandlePut)
	router.HandleFunc("DELETE /User/{id}", authHandler.HandleDelete)
	router.HandleFunc("POST /User/login", authHandler.HandleLogin)

	return router
}

func swaggerDocs() http.Handler {
	// Implement the handler for serving Swagger documentation
	return http.FileServer(http.Dir("./"))
}

func swaggerSpec(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./swagger.yaml")
}
