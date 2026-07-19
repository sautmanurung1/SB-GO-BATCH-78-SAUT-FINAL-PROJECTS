package services

import (
	"management-stock/models"
	"management-stock/repositories"
)

type MasterService interface {
	GetSectors() ([]models.Sector, error)
	CreateSector(sector *models.Sector) error
	UpdateSector(id int, sector *models.Sector) error
	DeleteSector(id int) error

	GetStocks() ([]models.Stock, error)
	CreateStock(stock *models.Stock) error
	UpdateStock(id int, stock *models.Stock) error
	DeleteStock(id int) error
}

type masterService struct {
	repo repositories.MasterRepository
}

func NewMasterService(repo repositories.MasterRepository) MasterService {
	return &masterService{repo}
}

func (s *masterService) GetSectors() ([]models.Sector, error) {
	return s.repo.GetSectors()
}
func (s *masterService) CreateSector(sector *models.Sector) error {
	return s.repo.CreateSector(sector)
}
func (s *masterService) UpdateSector(id int, sector *models.Sector) error {
	return s.repo.UpdateSector(id, sector)
}
func (s *masterService) DeleteSector(id int) error {
	return s.repo.DeleteSector(id)
}

func (s *masterService) GetStocks() ([]models.Stock, error) {
	return s.repo.GetStocks()
}
func (s *masterService) CreateStock(stock *models.Stock) error {
	return s.repo.CreateStock(stock)
}
func (s *masterService) UpdateStock(id int, stock *models.Stock) error {
	return s.repo.UpdateStock(id, stock)
}
func (s *masterService) DeleteStock(id int) error {
	return s.repo.DeleteStock(id)
}
