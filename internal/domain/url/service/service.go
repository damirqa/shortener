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

func (s *URLService) SaveURL(shortURL, longURL *entity.URL) {
	err := s.repo.Insert(shortURL.Link, *longURL)
	if err != nil {
		logger.GetLogger().Error("cannot insert link", zap.Error(err))
	}
}

func (s *URLService) GenerateShortURL() *entity.URL {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			logger.GetLogger().Error("Problem with generate short URL", zap.Error(err))
		}

		b[i] = letters[index.Int64()]
	}

	return entity.New(string(b))
}

func (s *URLService) Get(shortURL *entity.URL) (entity.URL, bool) {
	longURL, exist, err := s.repo.Get(shortURL.Link)
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
	for key, URLEntity := range urls {
		urlData := URLData{
			UUID:        key,
			ShortURL:    key,
			OriginalURL: URLEntity.Link,
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

		url := entity.New(urlData.OriginalURL)

		err = s.repo.Insert(urlData.UUID, *url)
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
		shortURL := s.GenerateShortURL()
		err := s.repo.InsertURLWithCorrelationID(shortURL.Link, url.OriginalURL)
		if err != nil {
			return nil, err
		}

		// (*shortURL).CorrelationID = url.CorrelationID
		if shortURL != nil {
			shortURL.CorrelationID = url.CorrelationID
		} else {
			logger.GetLogger().Error("DANGEROUS PLACE VSE TAKI NIL")
		}
		res = append(res, shortURL)
	}

	return res, nil
}
