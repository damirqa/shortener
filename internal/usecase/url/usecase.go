package url

import (
	"errors"
	"github.com/damirqa/shortener/cmd/config"
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

func (u UseCase) Generate(longURL string) (*URLDomainEntity.URL, error) {
	longURLEntity := URLDomainEntity.New(longURL)
	shortURLEntity := u.service.GenerateShortURL()
	err := u.service.SaveURL(shortURLEntity, longURLEntity)
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

	return shortURLEntity, nil
}

func (u UseCase) Get(shortURL string) (URLDomainEntity.URL, bool) {
	shortURLEntity := URLDomainEntity.New(shortURL)
	if shortURLEntity == nil {
		return URLDomainEntity.URL{}, false
	}

	longURL, exist := u.service.Get(shortURLEntity)
	return longURL, exist
}

func (u UseCase) GenerateBatch(request []model.URLRequestWithCorrelationID) ([]model.URLResponseWithCorrelationID, error) {
	URLSRequestWithCorrelationID := make([]model.URLRequestWithCorrelationID, 0, 1000)
	var URLSResponseWithCorrelationID []model.URLResponseWithCorrelationID
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

	for _, e := range URLEntities {
		URLResponseWithCorrelationID := model.URLResponseWithCorrelationID{CorrelationID: e.CorrelationID, ShortURL: config.Instance.GetResultAddress() + "/" + e.Link}
		URLSResponseWithCorrelationID = append(URLSResponseWithCorrelationID, URLResponseWithCorrelationID)
	}

	return URLSResponseWithCorrelationID, nil
}
