package models

import "time"

type Sector struct {
	ID         int       `json:"id" db:"id"`
	NamaSektor string    `json:"nama_sektor" db:"nama_sektor" binding:"required"`
	Deskripsi  string    `json:"deskripsi" db:"deskripsi"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type Stock struct {
	ID           int       `json:"id" db:"id"`
	TickerSymbol string    `json:"ticker_symbol" db:"ticker_symbol" binding:"required"`
	CompanyName  string    `json:"company_name" db:"company_name" binding:"required"`
	SectorID     int       `json:"sector_id" db:"sector_id" binding:"required"`
	CurrentPrice float64   `json:"current_price" db:"current_price" binding:"required"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
