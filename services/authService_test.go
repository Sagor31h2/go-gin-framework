package services

import (
	"gin-test/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{}, &models.Note{}, &models.NoteShare{}, &models.Attachment{})
	return db
}

func TestAuthService_Register(t *testing.T) {
	db := setupTestDB()
	s := NewAuthService(db)

	t.Run("Success", func(t *testing.T) {
		user, err := s.Register("test@example.com", "password123")
		assert.NoError(t, err)
		assert.NotZero(t, user.ID)
		assert.NotEmpty(t, user.PasswordHash)
	})

	t.Run("Duplicate Email", func(t *testing.T) {
		_, _ = s.Register("dup@example.com", "password123")
		_, err := s.Register("dup@example.com", "password456")
		assert.Error(t, err)
	})
}

func TestAuthService_Login(t *testing.T) {
	db := setupTestDB()
	s := NewAuthService(db)

	// Seed user
	_, _ = s.Register("login@example.com", "password123")

	t.Run("Success", func(t *testing.T) {
		token, err := s.Login("login@example.com", "password123")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Wrong Password", func(t *testing.T) {
		token, err := s.Login("login@example.com", "wrongpass")
		assert.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("User Not Found", func(t *testing.T) {
		token, err := s.Login("nonexistent@example.com", "password123")
		assert.Error(t, err)
		assert.Empty(t, token)
	})
}
