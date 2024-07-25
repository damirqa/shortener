package service

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"github.com/damirqa/shortener/cmd/config"
	"github.com/damirqa/shortener/internal/domain/url/entity"
	"github.com/damirqa/shortener/internal/domain/url/model"
	"github.com/damirqa/shortener/internal/domain/url/repository"
	"github.com/damirqa/shortener/internal/infrastructure/logger"
	"go.uber.org/zap"
	"math/big"
	"os"
)

type URLService struct {
	repo repository.URLRepository
}

func New(repo repository.URLRepository) *URLService {
	return &URLService{
		repo: repo,
	}
}

func (s *URLService) SaveURL(URLEntity *entity.URL) error {
	err := s.repo.Insert(URLEntity)
	if err != nil {
		return err
	}

	return nil
}

func (s *URLService) GenerateShortURL() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			logger.GetLogger().Error("Problem with generate short URL", zap.Error(err))
		}

		b[i] = letters[index.Int64()]
	}

	return string(b)
}

func (s *URLService) Get(URLEntity *entity.URL) (*entity.URL, bool) {
	longURL, exist, err := s.repo.Get(URLEntity.ShortURL)
	if err != nil {
		logger.GetLogger().Error("cannot get url from db", zap.Error(err))
	}

	return longURL, exist
}

type URLData struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (s *URLService) SaveToFile() {
	urls, err := s.repo.GetAll()
	if err != nil {
		logger.GetLogger().Error("problem get all data from urls", zap.Error(err))
	}

	file, err := os.OpenFile(config.Instance.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logger.GetLogger().Error("Error opening file", zap.Error(err))
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.GetLogger().Error(err.Error())
		}
	}(file)

	encoder := json.NewEncoder(file)
	for _, URLEntity := range urls {
		urlData := URLData{
			UUID:        URLEntity.ShortURL,
			ShortURL:    URLEntity.ShortURL,
			OriginalURL: URLEntity.OriginalURL,
			// todo: может добавить userID?
		}
		if err := encoder.Encode(&urlData); err != nil {
			logger.GetLogger().Error("Error encoding url data", zap.Error(err))
		}
	}
}

func (s *URLService) LoadFromFile() {
	if _, err := os.Stat(config.Instance.FileStoragePath); os.IsNotExist(err) {
		logger.GetLogger().Info("File not exist, skipping load:", zap.String("path", config.Instance.FileStoragePath))
		return
	}

	file, err := os.Open(config.Instance.FileStoragePath)
	if err != nil {
		logger.GetLogger().Error(err.Error())
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.GetLogger().Error(err.Error())
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue // Пропускаем пустые строки
		}

		var urlData URLData

		err := json.Unmarshal([]byte(line), &urlData)
		if err != nil {
			logger.GetLogger().Error(err.Error())
		}

		url := entity.URL{ShortURL: urlData.ShortURL, OriginalURL: urlData.OriginalURL}

		err = s.repo.Insert(&url)
		if err != nil {
			logger.GetLogger().Error("cannot insert link", zap.Error(err))
		}
	}

	if err := scanner.Err(); err != nil {
		logger.GetLogger().Error(err.Error())
	}
}

func (s *URLService) CreateURLs(urls []model.URLRequestWithCorrelationID) ([]*entity.URL, error) {
	res := make([]*entity.URL, 0, len(urls))

	for _, url := range urls {
		// todo: здесь нужно проверять на существующие ссылки, поломает ли это тесты?
		shortURL := s.GenerateShortURL()
		u := entity.URL{ShortURL: shortURL, OriginalURL: url.OriginalURL}
		err := s.repo.InsertURLWithCorrelationID(u.ShortURL, u.OriginalURL)
		if err != nil {
			return nil, err
		}

		u.CorrelationID = url.CorrelationID
		res = append(res, &u)
	}

	return res, nil
}

func (s *URLService) GetShortURLByOriginalURL(longURL string) (*entity.URL, error) {
	URLEntity, err := s.repo.FindByOriginalURL(longURL)
	if err != nil {
		return nil, err
	}

	return URLEntity, nil
}

func (s *URLService) GetAllUserLinks(userID string) ([]*entity.URL, error) {
	urls, err := s.repo.GetAllUserLinks(userID)
	if err != nil {
		return nil, err
	}

	return urls, nil
}
