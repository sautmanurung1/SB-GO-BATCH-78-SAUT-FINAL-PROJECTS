package models

import "time"

type Transaction struct {
	ID              string    `json:"id" db:"id"`
	UserID          string    `json:"user_id" db:"user_id"`
	StockID         int       `json:"stock_id" db:"stock_id" binding:"required"`
	TransactionType string    `json:"transaction_type" db:"transaction_type" binding:"required,oneof=buy sell"`
	LotAmount       int       `json:"lot_amount" db:"lot_amount" binding:"required,gt=0"`
	PricePerShare   float64   `json:"price_per_share" db:"price_per_share" binding:"required,gt=0"`
	TransactionDate time.Time `json:"transaction_date" db:"transaction_date"`
}

type PortfolioItem struct {
	TickerSymbol string  `json:"ticker_symbol"`
	CompanyName  string  `json:"company_name"`
	TotalLot     int     `json:"total_lot"`
	AveragePrice float64 `json:"average_price"`
	CurrentValue float64 `json:"current_value"`
}
