package url

import "github.com/damirqa/shortener/internal/domain/url/entity"

type ServiceInterface interface {
	Generate(longURL string) []byte
	Get(shortURL string) (entity.URL, bool)
}
