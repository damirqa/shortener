package repository

import "github.com/damirqa/shortener/internal/domain/url/entity"

type Repo interface {
	Insert(key string, value entity.URL)
	Get(key string) (entity.URL, bool)
}
