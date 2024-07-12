package service

import (
	"github.com/damirqa/shortener/internal/domain/url/entity"
	"github.com/damirqa/shortener/internal/domain/url/model"
)

type BaseDomainService interface {
	GenerateShortURL() *entity.URL
	SaveURL(shortURL, longURL *entity.URL)
	Get(shortURL *entity.URL) (entity.URL, bool)
	CreateURLs(urls []model.URLRequestWithCorrelationID) ([]*entity.URL, error)
}
