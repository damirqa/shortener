package entity

type URL struct {
	ShortURL      string
	OriginalURL   string
	UserID        string
	CorrelationID string
	IsDeleted     bool
}

func New(shortURL, longURL, userID string) *URL {
	return &URL{
		ShortURL:    shortURL,
		OriginalURL: longURL,
		UserID:      userID,
	}
}
