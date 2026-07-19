package repositories

import (
	"database/sql"
	"management-stock/models"
)

type TransactionRepository interface {
	CreateTransaction(tx *models.Transaction) error
	GetTransactionsByUserID(userID string) ([]models.Transaction, error)
	GetPortfolioByUserID(userID string) ([]models.PortfolioItem, error)
	GetTotalLotByStockAndUser(userID string, stockID int) (int, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) CreateTransaction(tx *models.Transaction) error {
	var query string
	var err error
	if !tx.TransactionDate.IsZero() {
		query = `INSERT INTO transactions (user_id, stock_id, transaction_type, lot_amount, price_per_share, transaction_date)
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
		err = r.db.QueryRow(query, tx.UserID, tx.StockID, tx.TransactionType, tx.LotAmount, tx.PricePerShare, tx.TransactionDate).Scan(&tx.ID)
	} else {
		query = `INSERT INTO transactions (user_id, stock_id, transaction_type, lot_amount, price_per_share)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id, transaction_date`
		err = r.db.QueryRow(query, tx.UserID, tx.StockID, tx.TransactionType, tx.LotAmount, tx.PricePerShare).Scan(&tx.ID, &tx.TransactionDate)
	}
	return err
}

func (r *transactionRepository) GetTransactionsByUserID(userID string) ([]models.Transaction, error) {
	query := `SELECT id, user_id, stock_id, transaction_type, lot_amount, price_per_share, transaction_date FROM transactions WHERE user_id = $1 ORDER BY transaction_date DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []models.Transaction
	for rows.Next() {
		var tx models.Transaction
		if err := rows.Scan(&tx.ID, &tx.UserID, &tx.StockID, &tx.TransactionType, &tx.LotAmount, &tx.PricePerShare, &tx.TransactionDate); err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}
	return txs, nil
}

func (r *transactionRepository) GetTotalLotByStockAndUser(userID string, stockID int) (int, error) {
	query := `
		SELECT
			COALESCE(SUM(CASE WHEN transaction_type = 'buy' THEN lot_amount ELSE 0 END), 0) -
			COALESCE(SUM(CASE WHEN transaction_type = 'sell' THEN lot_amount ELSE 0 END), 0) as total_lot
		FROM transactions
		WHERE user_id = $1 AND stock_id = $2
	`
	var totalLot int
	err := r.db.QueryRow(query, userID, stockID).Scan(&totalLot)
	return totalLot, err
}

func (r *transactionRepository) GetPortfolioByUserID(userID string) ([]models.PortfolioItem, error) {
	query := `
		SELECT 
			s.ticker_symbol, 
			s.company_name,
			SUM(CASE WHEN t.transaction_type = 'buy' THEN t.lot_amount ELSE -t.lot_amount END) AS total_lot,
			COALESCE(
				SUM(CASE WHEN t.transaction_type = 'buy' THEN t.lot_amount * t.price_per_share ELSE 0 END) / 
				NULLIF(SUM(CASE WHEN t.transaction_type = 'buy' THEN t.lot_amount ELSE 0 END), 0)
			, 0) AS average_price,
			s.current_price * SUM(CASE WHEN t.transaction_type = 'buy' THEN t.lot_amount ELSE -t.lot_amount END) * 100 AS current_value
		FROM transactions t
		JOIN stocks s ON t.stock_id = s.id
		WHERE t.user_id = $1
		GROUP BY s.id, s.ticker_symbol, s.company_name, s.current_price
		HAVING SUM(CASE WHEN t.transaction_type = 'buy' THEN t.lot_amount ELSE -t.lot_amount END) > 0
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.PortfolioItem
	for rows.Next() {
		var item models.PortfolioItem
		if err := rows.Scan(&item.TickerSymbol, &item.CompanyName, &item.TotalLot, &item.AveragePrice, &item.CurrentValue); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
