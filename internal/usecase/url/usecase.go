package url

import (
	"fmt"
	"github.com/damirqa/shortener/cmd/config"
	URLDomainEntity "github.com/damirqa/shortener/internal/domain/url/entity"
	"github.com/damirqa/shortener/internal/domain/url/model"
	URLDomainService "github.com/damirqa/shortener/internal/domain/url/service"
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

func (u UseCase) Generate(longURL string) []byte {
	longURLEntity := URLDomainEntity.New(longURL)
	shortURLEntity := u.service.GenerateShortURL()
	u.service.SaveURL(shortURLEntity, longURLEntity)

	fullURL := fmt.Append([]byte(config.Instance.GetResultAddress()+"/"), shortURLEntity.Link)

	return fullURL
}

func (u UseCase) Get(shortURL string) (URLDomainEntity.URL, bool) {
	shortURLEntity := URLDomainEntity.New(shortURL)
	longURL, exist := u.service.Get(shortURLEntity)
	return longURL, exist
}

func (u UseCase) GenerateBatch(request []model.URLRequestWithCorrelationID) ([]model.URLResponseWithCorrelationID, error) {
	URLSRequestWithCorrelationID := make([]model.URLRequestWithCorrelationID, 0, 1000)
	var URLSResponseWithCorrelationID []model.URLResponseWithCorrelationID
	var URLEntities []*URLDomainEntity.URL

	for _, URLReq := range request {
		URLSRequestWithCorrelationID = append(URLSRequestWithCorrelationID, URLReq)

		if len(URLSRequestWithCorrelationID) == 1000 {
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
		URLResponseWithCorrelationID := model.URLResponseWithCorrelationID{CorrelationID: e.CorrelationID, ShortURL: e.Link}
		URLSResponseWithCorrelationID = append(URLSResponseWithCorrelationID, URLResponseWithCorrelationID)
	}

	return URLSResponseWithCorrelationID, nil
}
