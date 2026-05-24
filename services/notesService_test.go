package services

import (
	"context"
	"gin-test/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotesService_CRUD(t *testing.T) {
	db := setupTestDB()
	s, _ := NewNotesService(db)
	ctx := context.Background()

	userID := uint(1)

	t.Run("Create Note", func(t *testing.T) {
		note := &models.Note{Title: "Test Note", UserID: userID}
		err := s.CreateNote(ctx, note)
		assert.NoError(t, err)
		assert.NotZero(t, note.ID)
	})

	t.Run("Get Notes", func(t *testing.T) {
		notes, err := s.GetNotes(ctx, userID, "", 10, 1)
		assert.NoError(t, err)
		assert.Len(t, notes, 1)
		assert.Equal(t, "Test Note", notes[0].Title)
	})

	t.Run("Get Note By ID", func(t *testing.T) {
		note, err := s.GetNoteByID(ctx, "1", userID)
		assert.NoError(t, err)
		assert.Equal(t, "Test Note", note.Title)
	})

	t.Run("Update Note", func(t *testing.T) {
		update := &models.Note{Title: "Updated Note"}
		err := s.UpdateNote(ctx, "1", userID, update)
		assert.NoError(t, err)

		note, _ := s.GetNoteByID(ctx, "1", userID)
		assert.Equal(t, "Updated Note", note.Title)
	})

	t.Run("Delete Note", func(t *testing.T) {
		err := s.DeleteNote(ctx, "1", userID)
		assert.NoError(t, err)

		_, err = s.GetNoteByID(ctx, "1", userID)
		assert.Error(t, err)
	})
}

func TestNotesService_Sharing(t *testing.T) {
	db := setupTestDB()
	s, _ := NewNotesService(db)
	ctx := context.Background()

	ownerID := uint(1)
	otherID := uint(2)

	note := &models.Note{Title: "Shared Note", UserID: ownerID}
	_ = s.CreateNote(ctx, note)
	noteID := uint(note.ID)

	t.Run("Share Note", func(t *testing.T) {
		err := s.ShareNote(ctx, noteID, ownerID, otherID)
		assert.NoError(t, err)

		shares, _ := s.ListShares(ctx, noteID, ownerID)
		assert.Len(t, shares, 1)
		assert.Equal(t, otherID, shares[0].SharedWithUser)
	})

	t.Run("Get Shared Note", func(t *testing.T) {
		notes, err := s.GetNotes(ctx, otherID, "", 10, 1)
		assert.NoError(t, err)
		assert.Len(t, notes, 1)
		assert.Equal(t, "Shared Note", notes[0].Title)
	})

	t.Run("Revoke Share", func(t *testing.T) {
		err := s.RevokeShare(ctx, noteID, ownerID, otherID)
		assert.NoError(t, err)

		shares, _ := s.ListShares(ctx, noteID, ownerID)
		assert.Len(t, shares, 0)
	})
}

func TestNotesService_Attachments(t *testing.T) {
	db := setupTestDB()
	s, _ := NewNotesService(db)
	ctx := context.Background()

	userID := uint(1)
	otherID := uint(2)

	note := &models.Note{Title: "Attachment Note", UserID: userID}
	_ = s.CreateNote(ctx, note)

	t.Run("Add Attachment Success", func(t *testing.T) {
		att, err := s.AddAttachment(ctx, note.ID, userID, "test.txt", "/path/to/test.txt")
		assert.NoError(t, err)
		assert.NotZero(t, att.ID)
		assert.Equal(t, "test.txt", att.FileName)
	})

	t.Run("Add Attachment Unauthorized", func(t *testing.T) {
		_, err := s.AddAttachment(ctx, note.ID, otherID, "hack.txt", "/path/to/hack.txt")
		assert.Error(t, err)
	})
}
