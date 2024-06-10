package url

import (
	"fmt"
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
	shortURL := u.service.GenerateShortURL()
	u.service.SaveURL(shortURL, longURL)

	fullURL := fmt.Append([]byte("http://127.0.0.1:8080/"), shortURL)

	return fullURL
}

func (u UseCase) Get(shortURL string) (string, bool) {
	longURL, exist := u.service.Get(shortURL)
	return longURL, exist
}
