package handlers

import (
	URLUseCase "github.com/damirqa/shortener/internal/usecase/url"
	"io"
	"log"
	"net/http"
	"strconv"
	"unicode/utf8"
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
				log.Fatalf("Error closing request body: %v", err)
			}
		}(r.Body)

		shortURL := useCase.Generate(string(longURL))

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", strconv.Itoa(utf8.RuneCount(shortURL)))
		w.WriteHeader(http.StatusCreated)

		_, err = w.Write(shortURL)
		if err != nil {
			http.Error(w, "Error generate short url", http.StatusInternalServerError)
			return
		}
	}
}
