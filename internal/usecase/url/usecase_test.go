package url

import (
	"bytes"
	urlDomainLocalRepository "github.com/damirqa/shortener/internal/domain/url/repository/local"
	urlDomainService "github.com/damirqa/shortener/internal/domain/url/service"

	"github.com/damirqa/shortener/internal/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerate(t *testing.T) {
	repo := urlDomainLocalRepository.New()
	service := urlDomainService.New(repo)
	useCase := New(service)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Shorten(useCase))

	server := httptest.NewServer(mux)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+"/", bytes.NewBufferString("https://practicum.yandex.ru"))
	if err != nil {
		t.Fatalf("Ошибка при попытке сделать запрос для сокращения URL. Ошибка: %v", err)
	}

	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusCreated {
		t.Errorf("Ожидался статус код %d, но получен %d", http.StatusCreated, status)
	}
}

func TestGet(t *testing.T) {
	repo := urlDomainLocalRepository.New()
	service := urlDomainService.New(repo)
	useCase := New(service)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Expand(useCase))

	server := httptest.NewServer(mux)
	defer server.Close()

	longURL := "http://detnkjoidndxr.ru/juc2om4xf"
	shortURL := service.GenerateShortURL()
	repo.Insert(shortURL, longURL)

	req, err := http.NewRequest(http.MethodGet, server.URL+"/"+shortURL, nil)
	if err != nil {
		t.Fatalf("Ошибка при попытке сделать запрос для получения полного URL. Ошибка: %v", err)
	}

	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("Ожидался статус код %d, но получен %d", http.StatusTemporaryRedirect, status)
	}

	if location := res.Header().Get("Location"); location != longURL {
		t.Errorf("Ожидался Location %s, но получен %s", longURL, location)
	}
}
