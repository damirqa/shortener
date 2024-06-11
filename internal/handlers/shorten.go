package handlers

import (
	URLUseCase "github.com/damirqa/shortener/internal/usecase/url"
	"io"
	"net/http"
	"strconv"
	"unicode/utf8"
)

func Shorten(useCase URLUseCase.ServiceInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		longURL, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		shortURL := useCase.Generate(string(longURL))

		w.Header().Set("Content-Type:", "text/plain")
		w.Header().Set("Content-Length", strconv.Itoa(utf8.RuneCount(shortURL)))
		w.WriteHeader(http.StatusCreated)

		_, err = w.Write(shortURL)
		if err != nil {
			http.Error(w, "Error generate short url", http.StatusInternalServerError)
			return
		}
	}
}
