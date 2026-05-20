package main

import (
	"log"

	"gin-test/controllers"
	"gin-test/internal/database"
	"gin-test/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Router
	router := gin.Default()

	// Initialize Database
	db, err := database.InitDb("test.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database connected successfully")

	// Initialize Services
	notesService, err := services.NewNotesService(db)
	if err != nil {
		log.Fatalf("Failed to initialize notes service: %v", err)
	}

	// Initialize Controllers
	controllers.NewNoteController(router, notesService)

	// Start Server
	log.Println("Starting server on :8000")
	if err := router.Run(":8000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
