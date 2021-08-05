package products

import "time"

type Model struct {
	ID           uint
	Name         string
	Price        uint
	Observations string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Storage interface {
	Migrate() error
}

type Service struct {
	storage Storage
}

func NewService(s Storage) *Service {
	return &Service{s}
}

func (s *Service) Migrate() error {
	return s.storage.Migrate()
}
