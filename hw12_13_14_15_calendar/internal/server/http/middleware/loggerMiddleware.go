package middleware

import (
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
)

type LoggerMiddleware struct {
	logger *zap.Logger
	next   http.Handler
}

func (l *LoggerMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	l.next.ServeHTTP(w, r)
	latency := time.Now().Sub(now)
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	l.logger.Debug(
		"new request",
		zap.String("ip", ip),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("proto", r.Proto),
		zap.Int64("latency", int64(latency.Seconds())),
		zap.String("user-agent", r.UserAgent()),
	)
}

func NewLoggerMiddleware(logger *zap.Logger, next http.Handler) http.Handler {
	return &LoggerMiddleware{logger, next}
}
