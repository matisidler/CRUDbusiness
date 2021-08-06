package sales

import (
	"fmt"
	"time"
)

type Model struct {
	ID           uint
	IdClient     uint
	IdProduct    uint
	Quantity     uint
	ProductPrice uint
	Income       int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Storage interface {
	Migrate() error
	Insert(*Model) error
	GetAll() ([]*Model, error)
	GetById(uint) (*Model, error)
	Update(*Model) error
	Delete(uint) error
}

type Service struct {
	storage Storage
}

func NewService(s Storage) *Service {
	return &Service{s}
}

func (m *Model) String() string {
	return fmt.Sprintf("%02d | %02d | %02d | %02d | %02d | %02d | %10s | %10s\n",
		m.ID, m.IdClient, m.IdProduct, m.Quantity, m.ProductPrice, m.Income, m.CreatedAt.Format("2006-01-02"), m.UpdatedAt.Format("2006-01-02"))
}

func (s *Service) Migrate() error {
	return s.storage.Migrate()
}

func (s *Service) Insert(m *Model) error {
	return s.storage.Insert(m)
}
func (s *Service) GetAll() ([]*Model, error) {
	return s.storage.GetAll()
}
func (s *Service) GetById(id uint) (*Model, error) {
	return s.storage.GetById(id)
}
func (s *Service) Update(m *Model) error {
	return s.storage.Update(m)
}
func (s *Service) Delete(id uint) error {
	return s.storage.Delete(id)
}
