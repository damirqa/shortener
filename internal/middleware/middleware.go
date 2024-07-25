package middleware

import (
	"compress/gzip"
	"context"
	"github.com/damirqa/shortener/cmd/config"
	"github.com/damirqa/shortener/internal/infrastructure/logger"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
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

// todo: почему лучше NewLoggingResponseWriterFunc, а не NewLoggingResponseWriter? я же возвращаю объект
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

// todo: стоит ли разные middleware выносить в разные файлы?
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

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

type contextUserKey string

const UserIDKey contextUserKey = "userID"

func CheckTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		cookie, err := request.Cookie("token")
		if err != nil {
			http.Error(writer, "Forbidden", http.StatusForbidden)
			return
		}

		tokenString := cookie.Value
		if tokenString == "" {
			http.Error(writer, "Forbidden", http.StatusForbidden)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Instance.SecretKey), nil
		}, jwt.WithoutClaimsValidation()) // todo: по идеи надо валидировать, чтобы выдать новый только если старый истек

		if err != nil {
			http.Error(writer, "Forbidden", http.StatusForbidden)
			return
		}

		if !token.Valid {
			http.Error(writer, "Forbidden", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(request.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func IssueTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authToken := request.Header.Get("Authorization")
		tokenString := ""
		if authToken != "" {
			tokenString = strings.TrimPrefix(authToken, "Bearer ")
		}

		userID := ""

		if tokenString == "" {
			userID = uuid.New().String()
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				},
				UserID: userID,
			})

			newTokenString, err := token.SignedString([]byte(config.Instance.SecretKey))
			if err != nil {
				http.Error(writer, "Internal server error", http.StatusInternalServerError)
				return
			}

			tokenString = newTokenString
		}

		cookie := http.Cookie{
			Name:  "token",
			Value: tokenString,
		}
		http.SetCookie(writer, &cookie)
		writer.Header().Set("Authorization", "Bearer "+tokenString)

		if userID == "" {
			claims := &Claims{}
			_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(config.Instance.SecretKey), nil
			}, jwt.WithoutClaimsValidation()) // todo: по идеи надо валидировать, чтобы выдать новый только если старый истек

			if err != nil {
				http.Error(writer, "Forbidden", http.StatusForbidden)
				return
			}

			userID = claims.UserID
		}

		ctx := context.WithValue(request.Context(), UserIDKey, userID)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
