package integration

import (
	"gin-test/internal/models"
	"gin-test/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}

	db.AutoMigrate(&models.User{}, &models.Note{}, &models.NoteShare{}, &models.Attachment{})
	return db
}

func SetupTestRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return routes.SetupRouter(db)
}
