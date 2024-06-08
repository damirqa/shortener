package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerate(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", generate)

	req, err := http.NewRequest(http.MethodPost, "/", bytes.NewBufferString("https://practicum.yandex.ru"))
	if err != nil {
		t.Fatalf("Ошибка при попытке сделать запрос для сокращения URL. Ошибка: %v", err)
	}

	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusCreated {
		t.Errorf("Ожидался статус код %d, но получен %d", http.StatusCreated, status)
	}
}

func TestGen(t *testing.T) {
	store = newURLStore()

	longURL := "http://detnkjoidndxr.ru/juc2om4xf"
	shortURL := generateShortURL()

	store.Set(shortURL, longURL)

	mux := http.NewServeMux()
	mux.HandleFunc("/", get)

	req, err := http.NewRequest(http.MethodGet, string(fmt.Append([]byte("/"), shortURL)), nil)
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
