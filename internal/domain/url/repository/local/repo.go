package local

import (
	"github.com/damirqa/shortener/internal/domain/url/entity"
	"sync"
)

// todo: в целом необходимо переписать
type URLLocalRepository struct {
	urls sync.Map
}

func (l *URLLocalRepository) Close() {

}

func (l *URLLocalRepository) InsertURLWithCorrelationID(short, long string) error {
	l.urls.Store(short, long)

	return nil
}

func New() *URLLocalRepository {
	return &URLLocalRepository{}
}

func (l *URLLocalRepository) Insert(URLEntity *entity.URL) error {
	l.urls.Store(URLEntity.ShortURL, URLEntity.OriginalURL)

	return nil
}

func (l *URLLocalRepository) Get(key string) (*entity.URL, bool, error) {
	value, exists := l.urls.Load(key)
	if !exists {
		return &entity.URL{}, false, nil
	}

	link, ok := value.(string)
	if !ok {
		return &entity.URL{}, false, nil
	}

	url := entity.URL{ShortURL: key, OriginalURL: link}

	return &url, true, nil
}

func (l *URLLocalRepository) GetAll() (map[string]*entity.URL, error) {
	urls := make(map[string]*entity.URL)
	l.urls.Range(func(key, value interface{}) bool {
		k, ok := key.(string)
		if !ok {
			return false
		}

		link, ok := value.(string)
		if !ok {
			return false
		}

		url := entity.URL{ShortURL: k, OriginalURL: link}
		urls[k] = &url

		return true
	})

	return urls, nil
}

func (l *URLLocalRepository) FindByOriginalURL(originalURL string) (*entity.URL, error) {
	//TODO implement me
	panic("implement me")
}

func (l *URLLocalRepository) GetAllUserLinks(userID string) ([]*entity.URL, error) {
	//TODO implement me
	panic("implement me")
}
