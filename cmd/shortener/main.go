package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type URLStore struct {
	sync.Mutex
	urls map[string]string
}

var store = newURLStore()

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	trimmedPath := strings.Trim(path, "/")

	splitPath := strings.Split(trimmedPath, "/")

	if len(splitPath) == 1 {
		runes := []rune(splitPath[0])
		count := len(runes)

		if count > 0 {
			if r.Method == http.MethodGet {
				get(w, r)
				return
			}

			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		} else {
			if r.Method == http.MethodPost {
				generate(w, r)
				return
			}

			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func generate(w http.ResponseWriter, r *http.Request) {
	link, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL := generateShortURL()
	store.Set(shortURL, string(link))

	fullURL := fmt.Append([]byte("http://127.0.0.1:8080/"), shortURL)

	w.Header().Set("Content-Type:", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(fullURL)))
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(fullURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	trimmedPath := strings.Trim(path, "/")

	longURL, exist := store.urls[trimmedPath]

	if !exist {
		http.Error(w, "URL not found", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
}

func newURLStore() *URLStore {
	return &URLStore{urls: make(map[string]string)}
}

func (s *URLStore) Set(shortURL, longURL string) {
	s.Lock()
	s.urls[shortURL] = longURL
	s.Unlock()
}

func (s *URLStore) Get(shortURL string) (string, bool) {
	url, exists := s.urls[shortURL]

	return url, exists
}

func generateShortURL() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
