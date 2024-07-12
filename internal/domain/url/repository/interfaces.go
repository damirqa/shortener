package repository

import "github.com/damirqa/shortener/internal/domain/url/entity"

type URLRepository interface {
	Insert(key string, value entity.URL) error
	Get(key string) (entity.URL, bool, error)
	GetAll() (map[string]entity.URL, error)
	InsertURLWithCorrelationID(short, long string) error
	Close()
}
