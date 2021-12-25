package middleware

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

type Logger interface {
	Write(msg string)
}

type LoggerMiddleware struct {
	logger Logger
	next   http.Handler
}

func (l *LoggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	l.next.ServeHTTP(w, r)
	latency := time.Now().Sub(now)
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	l.logger.Write(
		fmt.Sprintf(
			"%s [%s] %s %s %s %d %d %s",
			ip,
			now.Format("2006-01-02 15:04:05"),
			r.Method,
			r.URL.Path,
			r.Proto,
			200,
			int64(latency.Seconds()),
			r.UserAgent(),
		),
	)
}

func NewLoggerMiddleware(logger Logger, next http.Handler) http.Handler {
	return &LoggerMiddleware{logger, next}
}
