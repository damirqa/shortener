package handlers

import (
	"github.com/damirqa/shortener/internal/handlers/api"
	"github.com/damirqa/shortener/internal/usecase"
	"github.com/gorilla/mux"
)

func RegisterHandlers(router *mux.Router, useCases *usecase.UseCases) {
	// todo: добавить константы к путям
	router.HandleFunc("/", ShortenURL(useCases.URLUseCase)).Methods("POST")
	router.HandleFunc("/ping", Ping()).Methods("GET")
	router.HandleFunc("/{id}", ExpandURL(useCases.URLUseCase)).Methods("GET")
	router.HandleFunc("/api/shorten", api.ShortenURL(useCases.URLUseCase)).Methods("POST")
	router.HandleFunc("/api/shorten/batch", api.ShortenURLSBatch(useCases.URLUseCase)).Methods("POST")
}

func RegisterUserHandlers(router *mux.Router, useCases *usecase.UseCases) {
	router.HandleFunc("/urls", api.GetAllUserLinks(useCases.URLUseCase)).Methods("GET")
}
