package url_test

import (
	"bytes"
	URLDomainLocalRepository "github.com/damirqa/shortener/internal/domain/url/repository/local"
	URLDomainService "github.com/damirqa/shortener/internal/domain/url/service"
	URLUseCase "github.com/damirqa/shortener/internal/usecase/url"
	"net/http"

	"github.com/damirqa/shortener/internal/handlers"
	"github.com/gorilla/mux"
	"net/http/httptest"
	"testing"
)

func TestGenerate(t *testing.T) {
	repo := URLDomainLocalRepository.New()
	service := URLDomainService.New(repo)
	useCase := URLUseCase.New(service)

	router := mux.NewRouter()
	router.HandleFunc("/", handlers.Shorten(useCase)).Methods("POST")

	server := httptest.NewServer(router)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+"/", bytes.NewBufferString("https://practicum.yandex.ru"))
	if err != nil {
		t.Fatalf("Ошибка при попытке сделать запрос для сокращения URL. Ошибка: %v", err)
	}

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusCreated {
		t.Errorf("Ожидался статус код %d, но получен %d", http.StatusCreated, status)
	}
}

func TestGet(t *testing.T) {
	repo := URLDomainLocalRepository.New()
	service := URLDomainService.New(repo)
	useCase := URLUseCase.New(service)

	longURL := "http://detnkjoidndxr.ru/juc2om4xf"
	shortURL := service.GenerateShortURL()
	repo.Insert(shortURL, longURL)

	router := mux.NewRouter()
	router.HandleFunc("/{id}", handlers.Expand(useCase)).Methods("GET")

	server := httptest.NewServer(router)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL+"/"+shortURL, nil)
	if err != nil {
		t.Fatalf("Ошибка при попытке сделать запрос для получения полного URL. Ошибка: %v", err)
	}

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("Ожидался статус код %d, но получен %d", http.StatusTemporaryRedirect, status)
	}

	if location := res.Header().Get("Location"); location != longURL {
		t.Errorf("Ожидался Location %s, но получен %s", longURL, location)
	}
}
