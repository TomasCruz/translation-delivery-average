package service

type Service struct {
	db         Database
	windowSize int
}

func NewService(dBase Database, windowSize int) Service {
	return Service{db: dBase, windowSize: windowSize}
}
