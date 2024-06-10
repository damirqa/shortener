package local

import (
	"sync"
)

type Local struct {
	sync.Mutex
	urls map[string]string
}

func New() *Local {
	return &Local{urls: make(map[string]string)}
}

func (l *Local) Insert(key, value string) {
	l.Lock()
	l.urls[key] = value
	l.Unlock()
}

func (l *Local) Get(key string) (string, bool) {
	url, exists := l.urls[key]

	return url, exists
}
