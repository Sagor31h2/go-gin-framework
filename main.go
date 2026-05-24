package main

import (
	"log"
	"os/exec"
	"runtime"
	"time"

	"gin-test/internal/database"
	"gin-test/internal/models"
	"gin-test/routes"
)

// @title Gin Note App API
// @version 1.0
// @description This is a sample Gin framework server.
// @host localhost:8000
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		log.Printf("unsupported platform for auto-opening browser")
	}
	if err != nil {
		log.Printf("failed to open browser: %v", err)
	}
}

func main() {
	db, err := database.InitDb("test.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := db.AutoMigrate(&models.Note{}, &models.User{}, &models.NoteShare{}, &models.Attachment{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	router := routes.SetupRouter(db)

	go func() {
		// Wait a bit for server to start
		time.Sleep(1 * time.Second)
		openBrowser("http://localhost:8000/swagger/index.html")
	}()

	log.Println("Starting server on :8000")
	if err := router.Run(":8000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

