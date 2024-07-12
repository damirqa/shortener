package model

type URLRequestWithCorrelationId struct {
	CorrelationId string `json:"correlation_id"`
	OriginalUrl   string `json:"original_url"`
}

type URLResponseWithCorrelationId struct {
	CorrelationId string `json:"correlation_id"`
	ShortUrl      string `json:"short_url"`
}
