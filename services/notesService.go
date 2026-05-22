package services

import (
	"context"
	models "gin-test/internal/models"

	"gorm.io/gorm"
)

// NoteService defines the interface for note-related operations.
type NoteService interface {
	GetNotes(ctx context.Context, userID uint, search string, limit, page int) ([]models.Note, error)
	GetNoteByID(ctx context.Context, id string, userID uint) (*models.Note, error)
	CreateNote(ctx context.Context, note *models.Note) error
	UpdateNote(ctx context.Context, id string, userID uint, note *models.Note) error
	DeleteNote(ctx context.Context, id string, userID uint) error
}

type notesService struct {
	db *gorm.DB
}

func NewNotesService(db *gorm.DB) (NoteService, error) {
	return &notesService{db: db}, nil
}

func (s *notesService) GetNotes(ctx context.Context, userID uint, search string, limit, page int) ([]models.Note, error) {
	var notes []models.Note
	query := s.db.WithContext(ctx).Where("user_id = ?", userID)

	if search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}

	if limit > 0 {
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func (s *notesService) GetNoteByID(ctx context.Context, id string, userID uint) (*models.Note, error) {
	var note models.Note
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&note).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

func (s *notesService) CreateNote(ctx context.Context, note *models.Note) error {
	return s.db.WithContext(ctx).Create(note).Error
}

func (s *notesService) UpdateNote(ctx context.Context, id string, userID uint, note *models.Note) error {
	return s.db.WithContext(ctx).Model(&models.Note{}).Where("id = ? AND user_id = ?", id, userID).Updates(note).Error
}

func (s *notesService) DeleteNote(ctx context.Context, id string, userID uint) error {
	return s.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&models.Note{}).Error
}
