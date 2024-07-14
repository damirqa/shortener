package url

import (
	"github.com/damirqa/shortener/internal/domain/url/entity"
)
import "github.com/damirqa/shortener/internal/domain/url/model"

type UseCaseInterface interface {
	Generate(longURL string) (*entity.URL, error)
	Get(shortURL string) (entity.URL, bool)
	GenerateBatch(request []model.URLRequestWithCorrelationID) ([]model.URLResponseWithCorrelationID, error)
}
