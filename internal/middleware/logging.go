package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// LoggingMiddleware logs HTTP requests and responses.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap ResponseWriter to capture status code
		rw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Call the next handler
		next.ServeHTTP(rw, r)

		// Log method, path, status code, and duration
		log.Printf(
			"%s %s %d %s",
			r.Method,
			r.URL.String(),
			rw.Status(),
			time.Since(start),
		)
	})
}
