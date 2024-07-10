package handlers

import (
	"github.com/damirqa/shortener/internal/infrastructure/logger"
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

		shortURL := useCase.Generate(string(longURL))

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)

		_, err = w.Write(shortURL)
		if err != nil {
			http.Error(w, "Error generate short url", http.StatusInternalServerError)
			return
		}
	}
}
