package service

import (
	"github.com/damirqa/shortener/internal/domain/url/repository"
	"math/rand"
)

type Service struct {
	repo repository.Repo
}

func New(repo repository.Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) SaveURL(shortURL, longURL string) {
	s.repo.Insert(shortURL, longURL)
}

func (s *Service) GenerateShortURL() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *Service) Get(shortURL string) (string, bool) {
	longURL, exist := s.repo.Get(shortURL)

	return longURL, exist
}
