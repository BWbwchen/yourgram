package service

type DB interface {
	Query(string) ([]string, error)
}
