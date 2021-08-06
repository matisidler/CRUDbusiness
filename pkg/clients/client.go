package clients

import (
	"fmt"
	"time"
)

type Model struct {
	ID         uint
	Name       string
	Country    string
	PostalCode string
	Comment    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
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
	return fmt.Sprintf("%02d | %-20s | %-20s | %-20s | %-20s | %10s | %10s\n",
		m.ID, m.Name, m.Country, m.PostalCode, m.Comment, m.CreatedAt.Format("2006-01-02"), m.UpdatedAt.Format("2006-01-02"))
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
func (s *Service) getById(id uint) (*Model, error) {
	return s.storage.GetById(id)
}
func (s *Service) Update(m *Model) error {
	return s.storage.Update(m)
}
func (s *Service) Delete(id uint) error {
	return s.storage.Delete(id)
}
