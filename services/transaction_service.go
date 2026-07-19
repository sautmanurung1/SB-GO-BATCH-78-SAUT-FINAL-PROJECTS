package services

import (
	"bytes"
	"errors"
	"fmt"
	"management-stock/models"
	"management-stock/repositories"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ledongthuc/pdf"
)

type TransactionService interface {
	CreateTransaction(userID string, tx models.Transaction) error
	GetTransactionsByUserID(userID string) ([]models.Transaction, error)
	GetPortfolioByUserID(userID string) ([]models.PortfolioItem, error)
	ProcessPDFTransaction(userID string, fileData []byte) error
}

type transactionService struct {
	repo       repositories.TransactionRepository
	masterRepo repositories.MasterRepository
}

func NewTransactionService(repo repositories.TransactionRepository, masterRepo repositories.MasterRepository) TransactionService {
	return &transactionService{repo, masterRepo}
}

func (s *transactionService) CreateTransaction(userID string, tx models.Transaction) error {
	tx.UserID = userID

	if tx.TransactionType == "sell" {
		totalLot, err := s.repo.GetTotalLotByStockAndUser(userID, tx.StockID)
		if err != nil {
			return err
		}
		if tx.LotAmount > totalLot {
			return errors.New("insufficient lot amount for sell transaction")
		}
	}

	return s.repo.CreateTransaction(&tx)
}

func (s *transactionService) GetTransactionsByUserID(userID string) ([]models.Transaction, error) {
	return s.repo.GetTransactionsByUserID(userID)
}

func (s *transactionService) GetPortfolioByUserID(userID string) ([]models.PortfolioItem, error) {
	return s.repo.GetPortfolioByUserID(userID)
}

func (s *transactionService) ProcessPDFTransaction(userID string, fileData []byte) error {
	tmpFile, err := os.CreateTemp("", "trade-*.pdf")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if _, err := tmpFile.Write(fileData); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	stat, err := os.Stat(tmpPath)
	if err != nil {
		return err
	}

	f, err := os.Open(tmpPath)
	if err != nil {
		return err
	}
	defer f.Close()

	r, err := pdf.NewReader(f, stat.Size())
	if err != nil {
		return fmt.Errorf("failed to read pdf: %v", err)
	}

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return fmt.Errorf("failed to extract text from pdf: %v", err)
	}
	buf.ReadFrom(b)
	content := buf.String()

	// Ekstraksi sesuai PRD: Ticker, Type, Lot, Harga, dan Tanggal.
	reTicker := regexp.MustCompile(`(?i)Ticker:\s*([A-Z]{4})`)
	reType := regexp.MustCompile(`(?i)Type:\s*(BUY|SELL)`)
	reLot := regexp.MustCompile(`(?i)Lot:\s*(\d+)`)
	rePrice := regexp.MustCompile(`(?i)Price:\s*([\d\.]+)`)
	reDate := regexp.MustCompile(`(?i)Date:\s*(\d{4}-\d{2}-\d{2})`)

	tickerMatch := reTicker.FindStringSubmatch(content)
	typeMatch := reType.FindStringSubmatch(content)
	lotMatch := reLot.FindStringSubmatch(content)
	priceMatch := rePrice.FindStringSubmatch(content)
	dateMatch := reDate.FindStringSubmatch(content)

	if len(tickerMatch) < 2 || len(typeMatch) < 2 || len(lotMatch) < 2 || len(priceMatch) < 2 {
		return errors.New("failed to parse necessary data from PDF")
	}

	ticker := strings.ToUpper(tickerMatch[1])
	txType := strings.ToLower(typeMatch[1])
	lotStr := lotMatch[1]
	priceStr := priceMatch[1]

	lotAmount, err := strconv.Atoi(lotStr)
	if err != nil {
		return fmt.Errorf("invalid lot amount: %v", err)
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return fmt.Errorf("invalid price: %v", err)
	}

	stock, err := s.masterRepo.GetStockByTicker(ticker)
	if err != nil {
		return fmt.Errorf("error fetching stock info: %v", err)
	}
	if stock == nil {
		return fmt.Errorf("stock ticker %s not found in master data", ticker)
	}

	tx := models.Transaction{
		StockID:         stock.ID,
		TransactionType: txType,
		LotAmount:       lotAmount,
		PricePerShare:   price,
	}

	if len(dateMatch) > 1 {
		parsedDate, err := time.Parse("2006-01-02", dateMatch[1])
		if err == nil {
			tx.TransactionDate = parsedDate
		}
	}

	return s.CreateTransaction(userID, tx)
}
