package entity

type URL struct {
	link string
}

func New(link string) *URL {
	return &URL{
		link: link,
	}
}

func (u *URL) GetLink() string {
	return u.link
}
