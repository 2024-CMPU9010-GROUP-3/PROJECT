package routes

import (
	"log"
	"net/http"
	"strings"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/2024-CMPU9010-GROUP-3/PROJECT/docs"
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
	Router.HandleFunc("/docs/", httpSwagger.WrapHandler)
}

func AddRoute(r route) {
	log.Printf("Add route: %v\n", r.route)
	v1Router.Handle(r.route, http.StripPrefix(strings.TrimSuffix(r.route, "/"), r.handler))
}
