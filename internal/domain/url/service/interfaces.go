package service

import "github.com/damirqa/shortener/internal/domain/url/entity"

type BaseDomainService interface {
	GenerateShortURL() *entity.URL
	SaveURL(shortURL, longURL *entity.URL)
	Get(shortURL *entity.URL) (entity.URL, bool)
}
