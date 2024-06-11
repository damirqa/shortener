package service

import (
	"github.com/damirqa/shortener/internal/domain/url/entity"
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

func (s *Service) SaveURL(shortURL, longURL *entity.URL) {
	s.repo.Insert(shortURL.GetLink(), *longURL)
}

func (s *Service) GenerateShortURL() *entity.URL {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return entity.New(string(b))
}

func (s *Service) Get(shortURL *entity.URL) (entity.URL, bool) {
	longURL, exist := s.repo.Get(shortURL.GetLink())

	return longURL, exist
}
