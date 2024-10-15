package routes

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
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

// CORS 处理程序
func enableCORS(next http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedOrigins([]string{"https://localhost:3000"}),                  // 允许的源
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}), // 允许的请求方法
		handlers.AllowCredentials(),                                                  // 允许凭证
	)(next)
}

func AddRoute(r route) {
	log.Printf("Add route: %v\n", r.route)
	v1Router.Handle(r.route, http.StripPrefix(strings.TrimSuffix(r.route, "/"), enableCORS(r.handler))) // 使用 CORS 处理程序
}
