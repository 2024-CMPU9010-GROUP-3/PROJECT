package main

import (
	"log"
	"net/http"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request!")
	w.Write([]byte("World!"))
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/hello", handleRequest)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Listening on port 8080")
	server.ListenAndServe()
}
