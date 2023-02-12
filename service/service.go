package service

type Service struct {
	db Database
}

func NewService(dBase Database) Service {
	return Service{db: dBase}
}
