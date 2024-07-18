package url_test

import (
	"bytes"
	"encoding/json"
	"github.com/damirqa/shortener/cmd/config"
	URLDomainEntity "github.com/damirqa/shortener/internal/domain/url/entity"
	URLDomainLocalRepository "github.com/damirqa/shortener/internal/domain/url/repository/local"
	URLDomainService "github.com/damirqa/shortener/internal/domain/url/service"
	"github.com/damirqa/shortener/internal/handlers"
	"github.com/damirqa/shortener/internal/handlers/api"
	URLUseCase "github.com/damirqa/shortener/internal/usecase/url"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMain(m *testing.M) {
	config.Init()
	m.Run()
}

// todo: переписать и обновить тесты с использованием testify
func TestGenerate(t *testing.T) {
	cfg := config.Init()

	repo := URLDomainLocalRepository.New()
	service := URLDomainService.New(repo)
	useCase := URLUseCase.New(service)

	router := mux.NewRouter()
	router.HandleFunc("/", handlers.ShortenURL(useCase)).Methods("POST")

	// todo: как-то проверить, что он закрылся
	server := httptest.NewServer(router)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, "http://"+cfg.GetAddress()+"/", bytes.NewBufferString("https://practicum.yandex.ru"))
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
	cfg := config.Init()

	repo := URLDomainLocalRepository.New()
	service := URLDomainService.New(repo)
	useCase := URLUseCase.New(service)

	longURL := URLDomainEntity.New("http://detnkjoidndxr.ru/juc2om4xf")
	shortURL := service.GenerateShortURL()
	repo.Insert(shortURL.Link, *longURL)

	router := mux.NewRouter()
	router.HandleFunc("/{id}", handlers.ExpandURL(useCase)).Methods("GET")

	// todo: как-то проверить, что он закрылся
	server := httptest.NewServer(router)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, "http://"+cfg.GetAddress()+"/"+shortURL.Link, nil)
	if err != nil {
		t.Fatalf("Ошибка при попытке сделать запрос для получения полного URL. Ошибка: %v", err)
	}

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("Ожидался статус код %d, но получен %d", http.StatusTemporaryRedirect, status)
	}

	if location := res.Header().Get("Location"); location != longURL.Link {
		t.Errorf("Ожидался Location %s, но получен %s", longURL, location)
	}
}

func TestShorten(t *testing.T) {
	cfg := config.Init()

	repo := URLDomainLocalRepository.New()
	service := URLDomainService.New(repo)
	useCase := URLUseCase.New(service)

	router := mux.NewRouter()
	router.HandleFunc("/", handlers.ShortenURL(useCase)).Methods("POST")

	// todo: как-то проверить, что он закрылся
	server := httptest.NewServer(router)
	defer server.Close()

	urlRequest := api.URLRequest{Link: "https://practicum.yandex.ru"}

	urlRequestMarshal, err := json.Marshal(&urlRequest)
	if err != nil {
		t.Fatalf("Ошибка при сериализации. Ошибка: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "http://"+cfg.GetAddress()+"/", bytes.NewBuffer(urlRequestMarshal))
	if err != nil {
		t.Fatalf("Ошибка при попытке сделать запрос для сокращения URL. Ошибка: %v", err)
	}

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusCreated {
		t.Errorf("Ожидался статус код %d, но получен %d", http.StatusCreated, status)
	}
}
