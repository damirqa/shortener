package handlers

import (
	"github.com/damirqa/shortener/internal/usecase"
	"github.com/gorilla/mux"
)

func RegisterHandlers(router *mux.Router, useCases *usecase.UseCases) {
	router.HandleFunc("/", ShortenURL(useCases.URLUseCase)).Methods("POST")
	router.HandleFunc("/{id}", ExpandURL(useCases.URLUseCase)).Methods("GET")
}
