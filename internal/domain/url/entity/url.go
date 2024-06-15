package entity

type URL struct {
	Link string
}

func New(link string) *URL {
	return &URL{
		Link: link,
	}
}
