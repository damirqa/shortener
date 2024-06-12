package url_test

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/damirqa/shortener/cmd/config"
	URLDomainEntity "github.com/damirqa/shortener/internal/domain/url/entity"
	URLDomainLocalRepository "github.com/damirqa/shortener/internal/domain/url/repository/local"
	URLDomainService "github.com/damirqa/shortener/internal/domain/url/service"
	URLUseCase "github.com/damirqa/shortener/internal/usecase/url"
	"net"
	"net/http"
	"time"

	"github.com/damirqa/shortener/internal/handlers"
	"github.com/gorilla/mux"
	"net/http/httptest"
	"testing"
)

func TestMain(m *testing.M) {
	config.Init()
	m.Run()
}

func TestGenerate(t *testing.T) {
	repo := URLDomainLocalRepository.New()
	service := URLDomainService.New(repo)
	useCase := URLUseCase.New(service)

	router := mux.NewRouter()
	router.HandleFunc("/", handlers.Shorten(useCase)).Methods("POST")

	server := http.Server{
		Addr:    config.Config.Address + ":" + config.Config.Port,
		Handler: router,
	}

	defer func(server *http.Server) {
		err := server.Close()
		if err != nil {

		}
	}(&server)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("Не удалось запустить сервер: %v", err)
			return
		}
	}()

	if err := waitForPort(config.Config.Address, config.Config.Port, 10*time.Second); err != nil {
		t.Fatalf("Не удалось дождаться пока порт %s станет доступен для запроса: %v", config.Config.Port, err)
	}

	req, err := http.NewRequest(http.MethodPost, "http://"+config.Config.Address+":"+config.Config.Port+"/", bytes.NewBufferString("https://practicum.yandex.ru"))
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

	longURL := URLDomainEntity.New("http://detnkjoidndxr.ru/juc2om4xf")
	shortURL := service.GenerateShortURL()
	repo.Insert(shortURL.GetLink(), *longURL)

	router := mux.NewRouter()
	router.HandleFunc("/{id}", handlers.Expand(useCase)).Methods("GET")

	server := http.Server{
		Addr:    config.Config.Address + ":" + config.Config.Port,
		Handler: router,
	}

	defer func(server *http.Server) {
		err := server.Close()
		if err != nil {

		}
	}(&server)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("Не удалось запустить сервер: %v", err)
			return
		}
	}()

	if err := waitForPort(config.Config.Address, config.Config.Port, 10*time.Second); err != nil {
		t.Fatalf("Не удалось дождаться пока порт %s станет доступен для запроса: %v", config.Config.Port, err)
	}

	req, err := http.NewRequest(http.MethodGet, "http://"+config.Config.Address+":"+config.Config.Port+"/"+shortURL.GetLink(), nil)
	if err != nil {
		t.Fatalf("Ошибка при попытке сделать запрос для получения полного URL. Ошибка: %v", err)
	}

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("Ожидался статус код %d, но получен %d", http.StatusTemporaryRedirect, status)
	}

	if location := res.Header().Get("Location"); location != longURL.GetLink() {
		t.Errorf("Ожидался Location %s, но получен %s", longURL, location)
	}
}

func waitForPort(host string, port string, timeout time.Duration) error {
	address := net.JoinHostPort(host, port)
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", address, 1*time.Second)
		if err == nil {
			err := conn.Close()
			if err != nil {
				return err
			}
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("не удалось дождаться пока порт %s станет доступен для запроса", port)
}
