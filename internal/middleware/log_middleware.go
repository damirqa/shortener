package middleware

import (
	"github.com/damirqa/shortener/internal/infrastructure/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func NewLoggingResponseWriterFunc(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK, 0}
}

func (r *LoggingResponseWriter) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *LoggingResponseWriter) Write(b []byte) (int, error) {
	bytesWritten, err := r.ResponseWriter.Write(b)
	r.size += bytesWritten
	return bytesWritten, err
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := NewLoggingResponseWriterFunc(w)

		logger.GetLogger().Info("HTTP request",
			zap.String("method", r.Method),
			zap.String("url_path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr))

		next.ServeHTTP(lw, r)

		duration := time.Since(start).Milliseconds()

		logger.GetLogger().Info("request",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.String("remote_address", r.RemoteAddr),
			zap.Int64("duration", duration),
			zap.Int("status_code", lw.statusCode),
			zap.Int("size", lw.size))
	})
}
