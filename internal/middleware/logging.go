package middleware

import (
	"net/http"
	"time"

	"github.com/ndavidson/ingestion/pkg/logger"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
	body        []byte
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func (rw *responseWriter) Write(body []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	rw.body = append(rw.body, body...)
	return rw.ResponseWriter.Write(body)
}

// Logging returns a middleware that logs HTTP requests and responses
func Logging(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := wrapResponseWriter(w)

			log.Info("Request started",
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			)

			next.ServeHTTP(wrapped, r)

			log.Info("Request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"status", wrapped.Status(),
				"duration", time.Since(start),
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			)
		})
	}
}
