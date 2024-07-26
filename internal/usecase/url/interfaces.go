package url

import (
	"github.com/damirqa/shortener/internal/domain/url/entity"
)
import "github.com/damirqa/shortener/internal/domain/url/model"

type UseCaseInterface interface {
	Generate(longURL string, userID string) (*entity.URL, error)
	Get(shortURL, userID string) (*entity.URL, bool)
	GenerateBatch(request []model.URLRequestWithCorrelationID) ([]*entity.URL, error)
	GetAllUserLinks(userID string) ([]*entity.URL, error)
	DeleteUserLinks(userID string, shortURLs []string) error
}
