package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	
	"github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/middleware"
	"github.com/2024-CMPU9010-GROUP-3/PROJECT/internal/routes"
)

const defaultPort = 8080

func main() {
	var port int
	flag.IntVar(&port, "port", defaultPort, "the port the server listens on")
	flag.IntVar(&port, "p", defaultPort, "the port the server listens on (shorthand)")

	flag.Parse()

	router := routes.Router

	middlewares := middleware.CreateMiddlewareStack(
		middleware.Logging,
	)

	server := http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: middlewares(router),
	}

	log.Printf("Listening on port %v\n", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server: %v\n", err)
	}
}
