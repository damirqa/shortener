package repository

import "github.com/damirqa/shortener/internal/domain/url/entity"

type URLRepository interface {
	Insert(*entity.URL) error
	Get(key string) (*entity.URL, bool, error)
	GetAll() (map[string]*entity.URL, error)
	InsertURLWithCorrelationID(short, long string) error
	Close()
	FindByOriginalURL(originalURL string) (*entity.URL, error)
	GetAllUserLinks(userID string) ([]*entity.URL, error)
}
