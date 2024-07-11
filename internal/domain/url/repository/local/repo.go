package local

import (
	"github.com/damirqa/shortener/internal/domain/url/entity"
	"sync"
)

type URLLocalRepository struct {
	urls sync.Map
}

func New() *URLLocalRepository {
	return &URLLocalRepository{}
}

func (l *URLLocalRepository) Insert(key string, value entity.URL) error {
	l.urls.Store(key, value.Link)

	return nil
}

func (l *URLLocalRepository) Get(key string) (entity.URL, bool, error) {
	value, exists := l.urls.Load(key)
	if !exists {
		return entity.URL{}, false, nil
	}

	link, ok := value.(string)
	if !ok {
		return entity.URL{}, false, nil
	}

	url := entity.New(link)

	return *url, true, nil
}

func (l *URLLocalRepository) GetAll() (map[string]entity.URL, error) {
	urls := make(map[string]entity.URL)
	l.urls.Range(func(key, value interface{}) bool {
		k, ok := key.(string)
		if !ok {
			return false
		}

		link, ok := value.(string)
		if !ok {
			return false
		}

		url := entity.New(link)
		urls[k] = *url

		return true
	})

	return urls, nil
}
