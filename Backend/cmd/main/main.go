package main

import (
	"log"
	"net/http"

	"github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/middleware"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("World!"))
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/hello", handleRequest)

	middlewares := middleware.CreateMiddlewareStack(
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8080",
		Handler: middlewares(router),
	}

	log.Println("Listening on port 8080")
	server.ListenAndServe()
}
