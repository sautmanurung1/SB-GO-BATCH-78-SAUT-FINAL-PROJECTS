package controllers

import (
	"bytes"
	"io"
	"management-stock/models"
	"management-stock/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	service services.TransactionService
}

func NewTransactionController(service services.TransactionService) *TransactionController {
	return &TransactionController{service}
}

func (c *TransactionController) CreateTransaction(ctx *gin.Context) {
	userID := ctx.GetString("userID")

	var req models.Transaction
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.CreateTransaction(userID, req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Transaction created successfully"})
}

func (c *TransactionController) GetTransactions(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	txs, err := c.service.GetTransactionsByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, txs)
}

func (c *TransactionController) GetPortfolio(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	portfolio, err := c.service.GetPortfolioByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, portfolio)
}

func (c *TransactionController) ImportPDF(ctx *gin.Context) {
	userID := ctx.GetString("userID")

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	if header.Size > 5*1024*1024 {
		ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "File size exceeds 5MB limit"})
		return
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	
	mimeType := http.DetectContentType(buf.Bytes())
	if mimeType != "application/pdf" {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid file format, must be PDF"})
		return
	}

	if err := c.service.ProcessPDFTransaction(userID, buf.Bytes()); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "PDF imported and transaction recorded successfully"})
}
