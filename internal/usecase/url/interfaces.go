package url

type ServiceInterface interface {
	Generate(longURL string) []byte
	Get(shortURL string) (string, bool)
}
