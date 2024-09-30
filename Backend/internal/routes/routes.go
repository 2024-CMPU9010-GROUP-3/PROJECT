package routes

import (
	"log"
	"net/http"
	"strings"
)

type route struct {
	route   string
	handler *http.ServeMux
}

var Router = http.NewServeMux()
var v1Router = http.NewServeMux()

func init() {
	log.Println("init routes")
	Router.Handle("/v1/", http.StripPrefix("/v1", v1Router))
}

func AddRoute(r route) {
	log.Printf("Add route: %v\n", r.route)
	v1Router.Handle(r.route, http.StripPrefix(strings.TrimSuffix(r.route, "/"), r.handler))
}
