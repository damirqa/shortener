package url

import "github.com/damirqa/shortener/internal/domain/url/entity"

type UseCaseInterface interface {
	Generate(longURL string) []byte
	Get(shortURL string) (entity.URL, bool)
}
