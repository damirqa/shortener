package url

import (
	"errors"
	URLDomainEntity "github.com/damirqa/shortener/internal/domain/url/entity"
	"github.com/damirqa/shortener/internal/domain/url/model"
	URLDomainService "github.com/damirqa/shortener/internal/domain/url/service"
	dberror "github.com/damirqa/shortener/internal/error"
	"github.com/damirqa/shortener/internal/infrastructure/logger"
)

type UseCase struct {
	service URLDomainService.BaseDomainService
}

func New(service URLDomainService.BaseDomainService) *UseCase {
	return &UseCase{
		service: service,
	}
}

func (u UseCase) Generate(longURL, userID string) (*URLDomainEntity.URL, error) {
	shortURL := u.service.GenerateShortURL()
	URLEntity := URLDomainEntity.New(shortURL, longURL, userID)

	err := u.service.SaveURL(URLEntity)
	if err != nil {
		var uniqueErr *dberror.UniqueConstraintError
		if errors.As(err, &uniqueErr) {
			uniqueError := err
			url, err := u.service.GetShortURLByOriginalURL(longURL)
			if err != nil {
				return nil, err
			}

			return url, uniqueError
		} else {
			return nil, err
		}
	}

	return URLEntity, nil
}

func (u UseCase) Get(shortURL, userID string) (*URLDomainEntity.URL, bool) {
	URLEntity := URLDomainEntity.URL{ShortURL: shortURL, UserID: userID}

	longURL, exist := u.service.Get(&URLEntity)
	return longURL, exist
}

func (u UseCase) GenerateBatch(request []model.URLRequestWithCorrelationID) ([]*URLDomainEntity.URL, error) {
	URLSRequestWithCorrelationID := make([]model.URLRequestWithCorrelationID, 0, 1000)
	var URLEntities []*URLDomainEntity.URL

	for _, URLReq := range request {
		URLSRequestWithCorrelationID = append(URLSRequestWithCorrelationID, URLReq)

		if len(URLSRequestWithCorrelationID) == 2 {
			entities, err := u.service.CreateURLs(URLSRequestWithCorrelationID)
			if err != nil {
				logger.GetLogger().Error(err.Error())
				return nil, err
			}

			URLEntities = append(URLEntities, entities...)
			URLSRequestWithCorrelationID = URLSRequestWithCorrelationID[:0]
		}
	}

	if len(URLSRequestWithCorrelationID) > 0 {
		entities, err := u.service.CreateURLs(URLSRequestWithCorrelationID)
		if err != nil {
			logger.GetLogger().Error(err.Error())
			return nil, err
		}
		URLEntities = append(URLEntities, entities...)
	}

	return URLEntities, nil
}

func (u UseCase) GetAllUserLinks(userID string) ([]*URLDomainEntity.URL, error) {
	urls, err := u.service.GetAllUserLinks(userID)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func (u UseCase) DeleteUserLinks(userID string, shortURLs []string) error {
	err := u.service.DeleteUserLinks(userID, shortURLs)
	if err != nil {
		return err
	}

	return nil
}
