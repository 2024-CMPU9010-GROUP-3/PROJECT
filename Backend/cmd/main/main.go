package main

import (
	"log"
	"net/http"

	"github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/middleware"
	"github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/routes"
)

func main() {
	router := routes.Router

	middlewares := middleware.CreateMiddlewareStack(
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8081",
		Handler: middlewares(router),
	}

	log.Println("Listening on port 8081")
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server: %v/n", err)
	}
}
