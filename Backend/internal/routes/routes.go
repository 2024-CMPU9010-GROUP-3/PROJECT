package routes

import (
	"log"
	"net/http"
	"strings"

	"github.com/go-openapi/runtime/middleware"
)

// @title Example API
// @version 1.0
// @description This is a sample server for a Go application.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /v1

type route struct {
	route   string
	handler *http.ServeMux
}

var Router = http.NewServeMux()
var v1Router = http.NewServeMux()

// @title Initialize Routes
// @description Initialize the main router and version 1 router
func init() {
	log.Println("init routes")

	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	Router.Handle("/docs", sh)
	Router.HandleFunc("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.yaml")
	})

	Router.Handle("/v1/", http.StripPrefix("/v1", v1Router))
}

// @Summary Add a new route
// @Description Add a new route to the version 1 router
// @Param route body route true "Route"
// @Success 200 {string} string "ok"
// @Router /addRoute [post]
func AddRoute(r route) {
	log.Printf("Add route: %v\n", r.route)
	v1Router.Handle(r.route, http.StripPrefix(strings.TrimSuffix(r.route, "/"), r.handler))
}
