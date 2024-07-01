package service

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"github.com/damirqa/shortener/cmd/config"
	"github.com/damirqa/shortener/internal/domain/url/entity"
	"github.com/damirqa/shortener/internal/domain/url/repository"
	"log"
	"math/big"
	"os"
)

type URLService struct {
	repo repository.URLRepository
}

func New(repo repository.URLRepository) *URLService {
	return &URLService{
		repo: repo,
	}
}

func (s *URLService) SaveURL(shortURL, longURL *entity.URL) {
	s.repo.Insert(shortURL.Link, *longURL)
}

func (s *URLService) GenerateShortURL() *entity.URL {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			log.Fatalf("Problem with generate short URL: %v", err)
		}

		b[i] = letters[index.Int64()]
	}

	return entity.New(string(b))
}

func (s *URLService) Get(shortURL *entity.URL) (entity.URL, bool) {
	longURL, exist := s.repo.Get(shortURL.Link)

	return longURL, exist
}

type URLData struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (s *URLService) DumpToFile() {
	urls := s.repo.GetAll()

	file, err := os.OpenFile(config.Instance.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	for key, URLEntity := range urls {
		urlData := URLData{
			UUID:        key,
			ShortURL:    key,
			OriginalURL: URLEntity.Link,
		}
		data, err := json.Marshal(urlData)
		if err != nil {
			log.Fatalf("error marshalling url data: %v", err)
		}

		_, err = file.Write(append(data, '\n'))
		if err != nil {
			log.Fatalf("error writing to file: %v", err)
		}
	}
}

func (s *URLService) RecoveryURLsFromFile() {
	if _, err := os.Stat(config.Instance.FileStoragePath); os.IsNotExist(err) {
		log.Println("File not exist, skipping load:", config.Instance.FileStoragePath)
		return
	}

	file, err := os.Open(config.Instance.FileStoragePath)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue // Пропускаем пустые строки
		}

		var urlData URLData
		err := json.Unmarshal([]byte(line), &urlData)
		if err != nil {
			log.Fatalf("error unmarshalling url data: %v", err)
		}
		url := entity.New(urlData.OriginalURL)
		s.repo.Insert(urlData.UUID, *url)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %v", err)
	}
}
