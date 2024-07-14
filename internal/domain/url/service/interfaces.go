package service

import (
	"github.com/damirqa/shortener/internal/domain/url/entity"
	"github.com/damirqa/shortener/internal/domain/url/model"
)

type BaseDomainService interface {
	GenerateShortURL() *entity.URL
	SaveURL(shortURL, longURL *entity.URL) error
	Get(shortURL *entity.URL) (entity.URL, bool)
	CreateURLs(urls []model.URLRequestWithCorrelationID) ([]*entity.URL, error)
	GetShortURLByOriginalURL(longURL string) (*entity.URL, error)
}
