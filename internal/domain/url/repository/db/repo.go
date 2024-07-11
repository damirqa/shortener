package db

import (
	"context"
	"errors"
	"github.com/damirqa/shortener/cmd/config"
	"github.com/damirqa/shortener/internal/domain/url/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type URLDBRepository struct {
	pool *pgxpool.Pool
}

func New() (*URLDBRepository, error) {
	pool, err := pgxpool.New(context.Background(), config.Instance.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	return &URLDBRepository{pool: pool}, nil
}

func (l *URLDBRepository) Insert(key string, value entity.URL) error {
	//l.urls.Store(key, value.Link)

	_, err := l.pool.Exec(context.Background(), "INSERT INTO urls (short, long) VALUES ($1, $2)", key, value.Link)
	return err
}

func (l *URLDBRepository) Get(key string) (entity.URL, bool, error) {
	//value, exists := l.urls.Load(key)
	//if !exists {
	//	return entity.URL{}, false
	//}
	//
	//link, ok := value.(string)
	//if !ok {
	//	return entity.URL{}, false
	//}
	//
	//url := entity.New(link)
	//
	//return *url, true

	var link string
	err := l.pool.QueryRow(context.Background(), "SELECT long FROM urls WHERE short = $1", key).Scan(&link)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.URL{}, false, nil
		}
		return entity.URL{}, false, err
	}

	url := entity.New(link)
	return *url, true, nil
}

func (l *URLDBRepository) GetAll() (map[string]entity.URL, error) {
	//urls := make(map[string]entity.URL)
	//l.urls.Range(func(key, value interface{}) bool {
	//	k, ok := key.(string)
	//	if !ok {
	//		return false
	//	}
	//
	//	link, ok := value.(string)
	//	if !ok {
	//		return false
	//	}
	//
	//	url := entity.New(link)
	//	urls[k] = *url
	//
	//	return true
	//})
	//
	//return urls

	rows, err := l.pool.Query(context.Background(), "SELECT short, long FROM urls")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urls := make(map[string]entity.URL)
	for rows.Next() {
		var key, link string
		if err := rows.Scan(&key, &link); err != nil {
			return nil, err
		}
		url := entity.New(link)
		urls[key] = *url
	}

	return urls, nil
}
