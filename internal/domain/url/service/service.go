package service

import (
	"crypto/rand"
	"github.com/damirqa/shortener/internal/domain/url/entity"
	"github.com/damirqa/shortener/internal/domain/url/repository"
	"log"
	"math/big"
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
