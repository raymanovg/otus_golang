package internalhttp

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/raymanovg/otus_golang/hw12_13_14_15_calendar/internal/app"
)

type responseWriter struct {
	http.ResponseWriter
	code int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func newLoggingMiddleware(logger app.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w}

			next.ServeHTTP(rw, r)

			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			logger.Info(fmt.Sprintf(
				"%s [%s] %s %s %s %d %d %s",
				ip,
				start.Format(time.RFC3339),
				r.Method,
				r.URL,
				r.Proto,
				rw.code,
				time.Since(start),
				r.UserAgent(),
			))
		})
	}
}
