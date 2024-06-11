package local

import (
	"github.com/damirqa/shortener/internal/domain/url/entity"
	"sync"
)

type Local struct {
	sync.Mutex
	urls map[string]entity.URL
}

func New() *Local {
	return &Local{urls: make(map[string]entity.URL)}
}

func (l *Local) Insert(key string, value entity.URL) {
	l.Lock()
	l.urls[key] = value
	l.Unlock()
}

func (l *Local) Get(key string) (entity.URL, bool) {
	url, exists := l.urls[key]

	return url, exists
}
