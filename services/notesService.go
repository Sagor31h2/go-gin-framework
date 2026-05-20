package services

import (
	"context"
	models "gin-test/internal/models"

	"gorm.io/gorm"
)

// NoteService defines the interface for note-related operations.
type NoteService interface {
	GetNotes(ctx context.Context) ([]models.Note, error)
	CreateNote(ctx context.Context, note *models.Note) error
}

type notesService struct {
	db *gorm.DB
}

// NewNotesService creates a new instance of NoteService.
func NewNotesService(db *gorm.DB) (NoteService, error) {
	if err := db.AutoMigrate(&models.Note{}); err != nil {
		return nil, err
	}
	return &notesService{db: db}, nil
}

func (s *notesService) GetNotes(ctx context.Context) ([]models.Note, error) {
	var notes []models.Note
	if err := s.db.WithContext(ctx).Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func (s *notesService) CreateNote(ctx context.Context, note *models.Note) error {
	return s.db.WithContext(ctx).Create(note).Error
}
