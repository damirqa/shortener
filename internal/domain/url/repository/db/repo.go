package db

import (
	"context"
	"errors"
	"github.com/damirqa/shortener/cmd/config"
	"github.com/damirqa/shortener/internal/domain/url/entity"
	"github.com/damirqa/shortener/internal/infrastructure/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type URLDBRepository struct {
	pool *pgxpool.Pool
}

func (l *URLDBRepository) Close() {
	l.pool.Close()
}

func New() (*URLDBRepository, error) {
	pool, err := pgxpool.New(context.Background(), config.Instance.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	err = conn.Conn().Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return &URLDBRepository{pool: pool}, nil
}

func (l *URLDBRepository) Insert(key string, value entity.URL) error {
	//l.urls.Store(key, value.Link)

	_, err := l.pool.Exec(context.Background(), "INSERT INTO urls (short, long) VALUES ($1, $2)", key, value.Link)
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}

	return nil
}

func (l *URLDBRepository) Get(key string) (entity.URL, bool, error) {
	var link string

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if l.pool == nil {
		logger.GetLogger().Fatal("pool is nil")
		return entity.URL{}, false, errors.New("pool is nil")
	}

	conn, err := l.pool.Acquire(ctx)
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}

	defer func() {
		conn.Release()
	}()

	_, err = conn.Conn().Prepare(ctx, "selectURL", "SELECT long FROM urls WHERE short = $1")
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}

	err = conn.Conn().QueryRow(ctx, "selectURL", key).Scan(&link)
	if err != nil {
		logger.GetLogger().Error(err.Error())

		return entity.URL{}, false, err
	}

	url := entity.New(link)
	return *url, true, nil
}

func (l *URLDBRepository) GetAll() (map[string]entity.URL, error) {
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

func (l *URLDBRepository) InsertURLWithCorrelationID(short, long string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if l.pool == nil {
		logger.GetLogger().Fatal("pool is nil")
		return errors.New("pool is nil")
	}

	conn, err := l.pool.Acquire(ctx)
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}

	defer func() {
		conn.Release()
	}()

	_, err = conn.Conn().Prepare(ctx, "insertURL", "INSERT INTO urls (short, long) VALUES ($1, $2)")
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}

	_, err = conn.Conn().Exec(ctx, "insertURL", short, long)
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}

	return nil
}
