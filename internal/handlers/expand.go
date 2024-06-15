package handlers

import (
	URLUseCase "github.com/damirqa/shortener/internal/usecase/url"
	"github.com/gorilla/mux"
	"net/http"
)

func ExpandURL(useCase URLUseCase.UseCaseInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		shortURL := vars["id"]

		longURL, exist := useCase.Get(shortURL)
		if !exist {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, longURL.Link, http.StatusTemporaryRedirect)
	}
}
