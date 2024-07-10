package middleware

import (
	"compress/gzip"
	"github.com/damirqa/shortener/internal/infrastructure/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
	"time"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

// todo: почему лучше NewLoggingResponseWriterFunc, а не NewLoggingResponseWriter
//
//	я же возвращаю объект
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

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func GzipMiddleware(next http.Handler) http.Handler {
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

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				logger.GetLogger().Fatal(err.Error())
			}
		}(request.Body)
	})
}
