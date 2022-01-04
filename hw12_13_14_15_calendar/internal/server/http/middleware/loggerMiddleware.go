package middleware

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type LoggerMiddleware struct {
	logger Logger
	next   http.Handler
}

func (l *LoggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	l.next.ServeHTTP(w, r)
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	latency := time.Now().Sub(now)

	l.logger.Debug(fmt.Sprintf(
		"ip: %s, method: %s, path: %s, proto: %s, latency: %d, user-agent: %s",
		ip, r.Method, r.URL.Path, r.Proto, latency, r.UserAgent(),
	))
}

func NewLoggerMiddleware(logger Logger, next http.Handler) http.Handler {
	return &LoggerMiddleware{logger, next}
}
