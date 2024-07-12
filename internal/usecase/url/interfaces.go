package url

import (
	"github.com/damirqa/shortener/internal/domain/url/entity"
)
import "github.com/damirqa/shortener/internal/domain/url/model"

type UseCaseInterface interface {
	Generate(longURL string) []byte
	Get(shortURL string) (entity.URL, bool)
	GenerateBatch(request []model.URLRequestWithCorrelationId) ([]model.URLResponseWithCorrelationId, error)
}
