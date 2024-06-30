package middleware

import (
	"compress/gzip"
	"github.com/damirqa/shortener/internal/infrastructure/logger"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
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

func LogMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := NewLoggingResponseWriter(w)

		log.Printf("%s - %s (%s)", r.Method, r.URL.Path, r.RemoteAddr)

		next.ServeHTTP(lw, r)

		duration := time.Since(start)

		logger.Log.Info("request",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.String("remote_address", r.RemoteAddr),
			zap.Int("duration", int(duration)),
			zap.Int("status_code", lw.statusCode),
			zap.Int("size", lw.size))
	})
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(request.Body)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			defer reader.Close()
			request.Body = io.NopCloser(reader)
		}

		if strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
			writer.Header().Set("Content-Encoding", "gzip")
			gzipWriter := gzip.NewWriter(writer)
			defer gzipWriter.Close()
			gzipResponseWriter := gzipResponseWriter{Writer: gzipWriter, ResponseWriter: writer}
			next.ServeHTTP(gzipResponseWriter, request)
		} else {
			next.ServeHTTP(writer, request)
		}
	})
}
