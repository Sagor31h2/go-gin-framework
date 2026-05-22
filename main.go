package main

import (
	"log"

	"gin-test/internal/database"
	"gin-test/internal/models"
	"gin-test/routes"
)

func main() {
	db, err := database.InitDb("test.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := db.AutoMigrate(&models.Note{}, &models.User{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	router := routes.SetupRouter(db)

	log.Println("Starting server on :8000")
	if err := router.Run(":8000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
