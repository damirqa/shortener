package url

import (
	"fmt"
	"github.com/damirqa/shortener/cmd/config"
	URLDomainEntity "github.com/damirqa/shortener/internal/domain/url/entity"
	URLDomainService "github.com/damirqa/shortener/internal/domain/url/service"
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

	fullURL := fmt.Append([]byte(config.ConfigInstance.GetResultAddress()+"/"), shortURLEntity.GetLink())

	return fullURL
}

func (u UseCase) Get(shortURL string) (URLDomainEntity.URL, bool) {
	shortURLEntity := URLDomainEntity.New(shortURL)
	longURL, exist := u.service.Get(shortURLEntity)
	return longURL, exist
}
