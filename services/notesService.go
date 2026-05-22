package services

import (
	"context"
	apperrors "gin-test/internal/errors"
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
	ShareNote(ctx context.Context, noteID uint, ownerID uint, targetUserID uint) error
	RevokeShare(ctx context.Context, noteID uint, ownerID uint, targetUserID uint) error
	ListShares(ctx context.Context, noteID uint, ownerID uint) ([]models.NoteShare, error)
}

type notesService struct {
	db *gorm.DB
}

func NewNotesService(db *gorm.DB) (NoteService, error) {
	return &notesService{db: db}, nil
}

func (s *notesService) GetNotes(ctx context.Context, userID uint, search string, limit, page int) ([]models.Note, error) {
	var notes []models.Note

	// owned + shared with me
	query := s.db.WithContext(ctx).
		Where("user_id = ? OR id IN (SELECT note_id FROM note_shares WHERE shared_with_user = ?)", userID, userID)

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
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.NotFound("note not found")
		}
		return nil, apperrors.Internal("failed to fetch note")
	}
	return &note, nil
}

func (s *notesService) CreateNote(ctx context.Context, note *models.Note) error {
	if err := s.db.WithContext(ctx).Create(note).Error; err != nil {
		return apperrors.Internal("failed to create note")
	}
	return nil
}

func (s *notesService) UpdateNote(ctx context.Context, id string, userID uint, note *models.Note) error {
	if err := s.db.WithContext(ctx).Model(&models.Note{}).Where("id = ? AND user_id = ?", id, userID).Updates(note).Error; err != nil {
		return apperrors.Internal("failed to update note")
	}
	return nil
}

func (s *notesService) DeleteNote(ctx context.Context, id string, userID uint) error {
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&models.Note{}).Error; err != nil {
		return apperrors.Internal("failed to delete note")
	}
	return nil
}

func (s *notesService) ShareNote(ctx context.Context, noteID uint, ownerID uint, targetUserID uint) error {
	var note models.Note
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", noteID, ownerID).First(&note).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return apperrors.NotFound("note not found")
		}
		return apperrors.Internal("failed to share note")
	}
	share := models.NoteShare{}
	return s.db.WithContext(ctx).
		Where(models.NoteShare{NoteID: noteID, OwnerID: ownerID, SharedWithUser: targetUserID}).
		Attrs(models.NoteShare{NoteID: noteID, OwnerID: ownerID, SharedWithUser: targetUserID}).
		FirstOrCreate(&share).Error
}

func (s *notesService) RevokeShare(ctx context.Context, noteID uint, ownerID uint, targetUserID uint) error {
	result := s.db.WithContext(ctx).
		Where("note_id = ? AND owner_id = ? AND shared_with_user = ?", noteID, ownerID, targetUserID).
		Delete(&models.NoteShare{})
	if result.Error != nil {
		return apperrors.Internal("failed to revoke share")
	}
	if result.RowsAffected == 0 {
		return apperrors.NotFound("share not found")
	}
	return nil
}

func (s *notesService) ListShares(ctx context.Context, noteID uint, ownerID uint) ([]models.NoteShare, error) {
	var shares []models.NoteShare
	if err := s.db.WithContext(ctx).
		Where("note_id = ? AND owner_id = ?", noteID, ownerID).
		Find(&shares).Error; err != nil {
		return nil, apperrors.Internal("failed to list shares")
	}
	return shares, nil
}
