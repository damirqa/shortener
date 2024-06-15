package local

import (
	"github.com/damirqa/shortener/internal/domain/url/entity"
	"sync"
)

type URLLocalRepository struct {
	sync.Mutex
	urls map[string]string
}

func New() *URLLocalRepository {
	return &URLLocalRepository{urls: make(map[string]string)}
}

func (l *URLLocalRepository) Insert(key string, value entity.URL) {
	l.Lock()
	defer l.Unlock()

	l.urls[key] = value.Link
}

func (l *URLLocalRepository) Get(key string) (entity.URL, bool) {
	link, exists := l.urls[key]
	url := entity.New(link)

	return *url, exists
}
