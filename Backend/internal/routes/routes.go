package routes

import (
	"log"
	"net/http"
	"strings"

	"github.com/2024-CMPU9010-GROUP-3/magpie/internal/handlers"
)

type route struct {
	route   string
	handler *http.ServeMux
}

var Router = http.NewServeMux()
var v1Router = http.NewServeMux()

func init() {
	Router.Handle("/v1/", http.StripPrefix("/v1", v1Router))
	Router.HandleFunc("GET /heartbeat", handlers.HandleHeartbeat)
}

func AddRoute(r route) {
	log.Printf("Add route: %v\n", r.route)
	v1Router.Handle(r.route, http.StripPrefix(strings.TrimSuffix(r.route, "/"), r.handler))
}
