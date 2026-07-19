package repositories

import (
	"database/sql"
	"management-stock/models"
)

type MasterRepository interface {
	GetSectors() ([]models.Sector, error)
	CreateSector(sector *models.Sector) error
	UpdateSector(id int, sector *models.Sector) error
	DeleteSector(id int) error

	GetStocks() ([]models.Stock, error)
	CreateStock(stock *models.Stock) error
	UpdateStock(id int, stock *models.Stock) error
	DeleteStock(id int) error
	GetStockByTicker(ticker string) (*models.Stock, error)
}

type masterRepository struct {
	db *sql.DB
}

func NewMasterRepository(db *sql.DB) MasterRepository {
	return &masterRepository{db}
}

func (r *masterRepository) GetSectors() ([]models.Sector, error) {
	rows, err := r.db.Query(`SELECT id, nama_sektor, deskripsi, created_at FROM stock_sectors`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sectors []models.Sector
	for rows.Next() {
		var s models.Sector
		if err := rows.Scan(&s.ID, &s.NamaSektor, &s.Deskripsi, &s.CreatedAt); err != nil {
			return nil, err
		}
		sectors = append(sectors, s)
	}
	return sectors, nil
}

func (r *masterRepository) CreateSector(sector *models.Sector) error {
	query := `INSERT INTO stock_sectors (nama_sektor, deskripsi) VALUES ($1, $2) RETURNING id, created_at`
	return r.db.QueryRow(query, sector.NamaSektor, sector.Deskripsi).Scan(&sector.ID, &sector.CreatedAt)
}

func (r *masterRepository) UpdateSector(id int, sector *models.Sector) error {
	query := `UPDATE stock_sectors SET nama_sektor = $1, deskripsi = $2 WHERE id = $3`
	_, err := r.db.Exec(query, sector.NamaSektor, sector.Deskripsi, id)
	return err
}

func (r *masterRepository) DeleteSector(id int) error {
	query := `DELETE FROM stock_sectors WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *masterRepository) GetStocks() ([]models.Stock, error) {
	rows, err := r.db.Query(`SELECT id, ticker_symbol, company_name, sector_id, current_price, created_at FROM stocks`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []models.Stock
	for rows.Next() {
		var s models.Stock
		if err := rows.Scan(&s.ID, &s.TickerSymbol, &s.CompanyName, &s.SectorID, &s.CurrentPrice, &s.CreatedAt); err != nil {
			return nil, err
		}
		stocks = append(stocks, s)
	}
	return stocks, nil
}

func (r *masterRepository) CreateStock(stock *models.Stock) error {
	query := `INSERT INTO stocks (ticker_symbol, company_name, sector_id, current_price) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	return r.db.QueryRow(query, stock.TickerSymbol, stock.CompanyName, stock.SectorID, stock.CurrentPrice).Scan(&stock.ID, &stock.CreatedAt)
}

func (r *masterRepository) UpdateStock(id int, stock *models.Stock) error {
	query := `UPDATE stocks SET ticker_symbol = $1, company_name = $2, sector_id = $3, current_price = $4 WHERE id = $5`
	_, err := r.db.Exec(query, stock.TickerSymbol, stock.CompanyName, stock.SectorID, stock.CurrentPrice, id)
	return err
}

func (r *masterRepository) DeleteStock(id int) error {
	query := `DELETE FROM stocks WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *masterRepository) GetStockByTicker(ticker string) (*models.Stock, error) {
	s := &models.Stock{}
	query := `SELECT id, ticker_symbol, company_name, sector_id, current_price, created_at FROM stocks WHERE ticker_symbol = $1`
	err := r.db.QueryRow(query, ticker).Scan(&s.ID, &s.TickerSymbol, &s.CompanyName, &s.SectorID, &s.CurrentPrice, &s.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s, nil
}
