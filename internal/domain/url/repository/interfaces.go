package repository

type Repo interface {
	Insert(key, value string)
	Get(key string) (string, bool)
}
