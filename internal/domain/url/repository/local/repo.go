package local

import (
	"github.com/damirqa/shortener/internal/domain/url/entity"
	"sync"
)

type URLLocalRepository struct {
	sync.Mutex
	urls map[string]entity.URL
}

func New() *URLLocalRepository {
	return &URLLocalRepository{urls: make(map[string]entity.URL)}
}

func (l *URLLocalRepository) Insert(key string, value entity.URL) {
	l.Lock()
	defer l.Unlock()

	l.urls[key] = value
}

func (l *URLLocalRepository) Get(key string) (entity.URL, bool) {
	url, exists := l.urls[key]

	return url, exists
}
