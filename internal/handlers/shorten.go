package handlers

import (
	"errors"
	"github.com/damirqa/shortener/cmd/config"
	dberror "github.com/damirqa/shortener/internal/error"
	"github.com/damirqa/shortener/internal/infrastructure/logger"
	"github.com/damirqa/shortener/internal/middleware"
	URLUseCase "github.com/damirqa/shortener/internal/usecase/url"
	"go.uber.org/zap"
	"io"
	"net/http"
)

func ShortenURL(useCase URLUseCase.UseCaseInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		longURL, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				logger.GetLogger().Error("Error closing request body", zap.Error(err))
			}
		}(r.Body)

		URLEntity, err := useCase.Generate(string(longURL), r.Context().Value(middleware.UserIDKey).(string))
		if err != nil {
			var uniqueErr *dberror.UniqueConstraintError
			if errors.As(err, &uniqueErr) {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusConflict)
			} else {
				logger.GetLogger().Error("Error generate short url", zap.Error(err))
				http.Error(w, "Error generate short url", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)

		fullURL := config.Instance.GetResultAddress() + "/" + URLEntity.ShortURL
		_, _ = w.Write([]byte(fullURL))
	}
}
