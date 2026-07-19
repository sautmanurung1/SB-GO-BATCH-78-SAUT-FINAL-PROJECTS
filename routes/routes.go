package routes

import (
	"database/sql"
	"management-stock/controllers"
	"management-stock/middlewares"
	"management-stock/repositories"
	"management-stock/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(db *sql.DB) *gin.Engine {
	r := gin.Default()

	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo)
	authController := controllers.NewAuthController(authService)

	masterRepo := repositories.NewMasterRepository(db)
	masterService := services.NewMasterService(masterRepo)
	masterController := controllers.NewMasterController(masterService)

	txRepo := repositories.NewTransactionRepository(db)
	txService := services.NewTransactionService(txRepo, masterRepo)
	txController := controllers.NewTransactionController(txService)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// Admin routes
		admin := api.Group("")
		admin.Use(middlewares.AuthMiddleware("admin"))
		{
			admin.GET("/sectors", masterController.GetSectors)
			admin.POST("/sectors", masterController.CreateSector)
			admin.PUT("/sectors/:id", masterController.UpdateSector)
			admin.DELETE("/sectors/:id", masterController.DeleteSector)

			admin.GET("/stocks", masterController.GetStocks)
			admin.POST("/stocks", masterController.CreateStock)
			admin.PUT("/stocks/:id", masterController.UpdateStock)
			admin.DELETE("/stocks/:id", masterController.DeleteStock)
		}

		// Investor routes
		investor := api.Group("")
		investor.Use(middlewares.AuthMiddleware("investor"))
		{
			investor.GET("/transactions", txController.GetTransactions)
			investor.POST("/transactions", txController.CreateTransaction)
			investor.POST("/transactions/import-pdf", txController.ImportPDF)
			investor.GET("/portfolio", txController.GetPortfolio)
		}
	}

	return r
}
