package main

import (
	"log"
	"management-stock/config"
	"management-stock/routes"
)

func main() {
	db := config.ConnectDB()
	defer db.Close()

	r := routes.SetupRoutes(db)

	log.Println("Server is running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
