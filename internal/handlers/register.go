package handlers

import (
	"github.com/damirqa/shortener/internal/usecase"
	"github.com/gorilla/mux"
)

func RegisterHandlers(router *mux.Router, useCases *usecase.UseCases) {
	router.HandleFunc("/", Shorten(useCases.URLUseCase)).Methods("POST")
	router.HandleFunc("/{id}", Expand(useCases.URLUseCase)).Methods("GET")
}
