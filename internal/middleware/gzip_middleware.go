package middleware

import (
	"compress/gzip"
	"github.com/damirqa/shortener/internal/infrastructure/logger"
	"io"
	"net/http"
	"strings"
)

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
