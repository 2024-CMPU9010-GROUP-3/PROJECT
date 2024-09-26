// This is a simple logging middleware for demonstration purposes
// created following this tutorial by DreamsOfCode on YouTube https://www.youtube.com/watch?v=H7tbjKFSg58
package middleware

import (
	"log"
	"net/http"
	"time"
)

type StatusCodeWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *StatusCodeWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		startTime := time.Now()

		statusCodeWriter := &StatusCodeWriter{
			ResponseWriter: writer,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(statusCodeWriter, request)

		log.Println(statusCodeWriter.statusCode, request.Method, request.URL.Path, time.Since(startTime))

	})
}
