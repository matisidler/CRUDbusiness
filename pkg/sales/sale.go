package sales

import "time"

type Model struct {
	ID        uint
	IdClient  uint
	IdProduct uint
	Quantity  uint
	Income    int
	CreatedAt time.Time
	UpdatedAt time.Time
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
