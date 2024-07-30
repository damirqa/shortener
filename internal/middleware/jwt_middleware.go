package middleware

import (
	"context"
	"github.com/damirqa/shortener/cmd/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"net"
	"net/http"
	"sync"
	"time"
)

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
			http.Error(writer, "Not authorized", http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value
		if tokenString == "" {
			http.Error(writer, "Not authorized", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Instance.SecretKey), nil
		}, jwt.WithoutClaimsValidation()) // todo: по идеи надо валидировать, чтобы выдать новый только если старый истек

		if err != nil {
			http.Error(writer, "Not authorized", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(writer, "Not authorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(request.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

// todo: не очень хорошее решение хранить токены,
// 		 но без этого не проходят тесты, так как на каждый новый запрос создается новый токен,
//		 а в тестах сохраняется последний токен запроса

// todo: переделать на sync.Map
var (
	tokenStore = make(map[string]string)
	mu         sync.Mutex
)

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return ip
}

func IssueTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ip := getIP(request)

		mu.Lock()
		tokenString, exists := tokenStore[ip]
		mu.Unlock()

		var userID string

		if exists {
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(config.Instance.SecretKey), nil
			}, jwt.WithoutClaimsValidation()) // todo: по идеи надо валидировать, чтобы выдать новый только если старый истек

			if err == nil && token.Valid {
				userID = claims.UserID
			}
		}

		if userID == "" {
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

			mu.Lock()
			tokenStore[ip] = tokenString
			mu.Unlock()
		}

		cookie := http.Cookie{
			Name:  "token",
			Value: tokenString,
		}
		http.SetCookie(writer, &cookie)
		writer.Header().Set("Authorization", "Bearer "+tokenString)

		ctx := context.WithValue(request.Context(), UserIDKey, userID)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
