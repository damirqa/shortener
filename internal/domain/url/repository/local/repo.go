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
	l.Lock()
	defer l.Unlock()

	link, exists := l.urls[key]
	url := entity.New(link)

	return *url, exists
}

func (l *URLLocalRepository) GetAll() map[string]entity.URL {
	l.Lock()
	defer l.Unlock()

	urls := make(map[string]entity.URL)
	for key, link := range l.urls {
		url := entity.New(link)
		urls[key] = *url
	}

	return urls
}
