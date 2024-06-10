package service

type BaseDomainService interface {
	GenerateShortURL() string
	SaveURL(shortURL, longURL string)
	Get(shortURL string) (string, bool)
}
