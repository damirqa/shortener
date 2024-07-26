package db

import (
	"context"
	"errors"
	"github.com/damirqa/shortener/cmd/config"
	"github.com/damirqa/shortener/internal/domain/url/entity"
	dberror "github.com/damirqa/shortener/internal/error"
	"github.com/damirqa/shortener/internal/infrastructure/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"sync"
)

type URLDBRepository struct {
	pool *pgxpool.Pool
}

// todo: нужно ли закрывать при завершении работы приложения?
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

func (l *URLDBRepository) Insert(URLEntity *entity.URL) error {
	_, err := l.pool.Exec(context.Background(), "INSERT INTO urls (short, long, user_id) VALUES ($1, $2, $3)", URLEntity.ShortURL, URLEntity.OriginalURL, URLEntity.UserID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return dberror.NewUniqueConstraintError(err)
		} else {
			return err
		}
	}

	return nil
}

func (l *URLDBRepository) Get(key string) (*entity.URL, bool, error) {
	var link string
	var is_deleted bool

	// todo: если использовать context.WithTimeout(context.Background(), 5*time.Second), то следующие запросы (в целом) не будут выполняться, почему?

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := l.pool.Acquire(ctx)
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}

	defer func() {
		conn.Release()
	}()

	_, err = conn.Conn().Prepare(ctx, "selectURL", "SELECT long, is_deleted FROM urls WHERE short = $1")
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}

	err = conn.Conn().QueryRow(ctx, "selectURL", key).Scan(&link, &is_deleted)
	if err != nil {
		logger.GetLogger().Error(err.Error())

		return &entity.URL{}, false, err
	}

	url := entity.URL{ShortURL: key, OriginalURL: link, IsDeleted: is_deleted}
	return &url, true, nil
}

func (l *URLDBRepository) GetAll() (map[string]*entity.URL, error) {
	rows, err := l.pool.Query(context.Background(), "SELECT short, long, user_id FROM urls")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urls := make(map[string]*entity.URL)
	for rows.Next() {
		var shortURL, originalURL, userID string
		if err := rows.Scan(&shortURL, &originalURL, &userID); err != nil {
			return nil, err
		}
		url := entity.New(shortURL, originalURL, userID)
		urls[shortURL] = url
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func (l *URLDBRepository) InsertURLWithCorrelationID(short, long string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

func (l *URLDBRepository) FindByOriginalURL(originalURL string) (*entity.URL, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := l.pool.Acquire(ctx)
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}

	defer conn.Release()

	_, err = conn.Conn().Prepare(ctx, "selectURL", "SELECT short FROM urls WHERE long = $1")
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}

	var short string
	err = conn.Conn().QueryRow(ctx, "selectURL", originalURL).Scan(&short)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return nil, err
	}

	return &entity.URL{ShortURL: short}, nil
}

func (l *URLDBRepository) GetAllUserLinks(userID string) ([]*entity.URL, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := l.pool.Acquire(ctx)
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}

	defer conn.Release()

	_, err = conn.Conn().Prepare(ctx, "selectURL", "SELECT short, long FROM urls WHERE user_id = $1")
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}

	rows, err := conn.Conn().Query(ctx, "selectURL", userID)
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}

	urls := make([]*entity.URL, 0)
	for rows.Next() {
		var short, long string
		if err := rows.Scan(&short, &long); err != nil {
			return nil, err
		}

		url := entity.New(config.Instance.GetResultAddress()+"/"+short, long, userID)
		urls = append(urls, url)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func (l *URLDBRepository) DeleteUserLinks(userID string, shortURLs []string) error {
	ErrorsCh := make(chan error, len(shortURLs))
	var wg sync.WaitGroup

	for _, url := range shortURLs {
		wg.Add(1)

		url := url
		go func() {
			defer wg.Done()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			conn, err := l.pool.Acquire(ctx)
			if err != nil {
				logger.GetLogger().Error(err.Error())
			}

			defer conn.Release()

			_, err = conn.Conn().Prepare(ctx, "deleteURL", "UPDATE urls SET is_deleted = true WHERE user_id = $1 AND short = $2")
			if err != nil {
				logger.GetLogger().Error(err.Error())
			}

			_, err = conn.Conn().Exec(ctx, "deleteURL", userID, url)
			if err != nil {
				logger.GetLogger().Error(err.Error())
				ErrorsCh <- err
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ErrorsCh)
	}()

	return nil
}
