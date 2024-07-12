package model

type URLRequestWithCorrelationID struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type URLResponseWithCorrelationId struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
