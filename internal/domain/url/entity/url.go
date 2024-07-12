package entity

type URL struct {
	Link          string
	CorrelationId string
}

func New(link string) *URL {
	return &URL{
		Link: link,
	}
}
