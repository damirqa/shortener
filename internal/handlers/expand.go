package handlers

import (
	"github.com/damirqa/shortener/internal/middleware"
	URLUseCase "github.com/damirqa/shortener/internal/usecase/url"
	"github.com/gorilla/mux"
	"net/http"
)

func ExpandURL(useCase URLUseCase.UseCaseInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		shortURL := vars["id"]

		URLEntity, exist := useCase.Get(shortURL, r.Context().Value(middleware.UserIDKey).(string))
		if !exist {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, URLEntity.OriginalURL, http.StatusTemporaryRedirect)
	}
}
