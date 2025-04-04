package cmd

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware is a middleware that logs request details.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		end := time.Now()
		duration := end.Sub(start)
		log.Printf("Response: %s %s %s %s", r.Method, r.URL.Path, http.StatusText(200), duration)
	})
}
