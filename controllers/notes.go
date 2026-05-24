package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	apperrors "gin-test/internal/errors"
	models "gin-test/internal/models"
	"gin-test/services"

	"github.com/gin-gonic/gin"
)

type NoteController struct {
	noteService services.NoteService
}

func NewNoteController(noteService services.NoteService) *NoteController {
	return &NoteController{noteService: noteService}
}

// GetNotes godoc
// @Summary List all notes
// @Description Get all owned and shared notes for the authenticated user
// @Tags notes
// @Security BearerAuth
// @Produce json
// @Param search query string false "Search by title"
// @Param limit query int false "Pagination limit" default(10)
// @Param page query int false "Pagination page" default(1)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /notes [get]
func (c *NoteController) GetNotes(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(uint)

	search := ctx.Query("search")
	limitStr := ctx.DefaultQuery("limit", "10")
	pageStr := ctx.DefaultQuery("page", "1")

	var limit, page int
	fmt.Sscanf(limitStr, "%d", &limit)
	fmt.Sscanf(pageStr, "%d", &page)
	if page < 1 {
		page = 1
	}

	notes, err := c.noteService.GetNotes(ctx.Request.Context(), userID, search, limit, page)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": notes, "limit": limit, "page": page})
}

// GetNoteByID godoc
// @Summary Get a note by ID
// @Description Fetch a specific note owned by the user
// @Tags notes
// @Security BearerAuth
// @Produce json
// @Param id path string true "Note ID"
// @Success 200 {object} map[string]models.Note
// @Failure 404 {object} map[string]string
// @Router /notes/{id} [get]
func (c *NoteController) GetNoteByID(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(uint)
	note, err := c.noteService.GetNoteByID(ctx.Request.Context(), ctx.Param("id"), userID)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": note})
}

// CreateNote godoc
// @Summary Create a new note
// @Description Add a new note to the user's collection
// @Tags notes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param note body models.Note true "Note object"
// @Success 201 {object} map[string]models.Note
// @Failure 400 {object} map[string]string
// @Router /notes [post]
func (c *NoteController) CreateNote(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(uint)
	var note models.Note
	if err := ctx.ShouldBindJSON(&note); err != nil {
		_ = ctx.Error(apperrors.BadRequest(err.Error()))
		ctx.Abort()
		return
	}
	note.UserID = userID
	if err := c.noteService.CreateNote(ctx.Request.Context(), &note); err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": note})
}

// UpdateNote godoc
// @Summary Update an existing note
// @Description Modify a note's title or status
// @Tags notes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Note ID"
// @Param note body models.Note true "Updated note object"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /notes/{id} [put]
func (c *NoteController) UpdateNote(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(uint)
	var note models.Note
	if err := ctx.ShouldBindJSON(&note); err != nil {
		_ = ctx.Error(apperrors.BadRequest(err.Error()))
		ctx.Abort()
		return
	}
	if err := c.noteService.UpdateNote(ctx.Request.Context(), ctx.Param("id"), userID, &note); err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "note updated successfully"})
}

// DeleteNote godoc
// @Summary Delete a note
// @Description Remove a note from the user's collection
// @Tags notes
// @Security BearerAuth
// @Param id path string true "Note ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /notes/{id} [delete]
func (c *NoteController) DeleteNote(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(uint)
	if err := c.noteService.DeleteNote(ctx.Request.Context(), ctx.Param("id"), userID); err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "note deleted successfully"})
}

// UploadAttachment godoc
// @Summary Upload an attachment to a note
// @Description Attach a file to a specific note
// @Tags notes
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Note ID"
// @Param file formData file true "File to upload"
// @Success 201 {object} map[string]models.Attachment
// @Failure 400 {object} map[string]string
// @Router /notes/{id}/attachments [post]
func (c *NoteController) UploadAttachment(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(uint)
	noteIDStr := ctx.Param("id")
	noteID, err := strconv.Atoi(noteIDStr)
	if err != nil {
		_ = ctx.Error(apperrors.BadRequest("invalid note id"))
		ctx.Abort()
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		_ = ctx.Error(apperrors.BadRequest("file is required"))
		ctx.Abort()
		return
	}

	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		_ = ctx.Error(apperrors.Internal("failed to create upload directory"))
		ctx.Abort()
		return
	}

	// Simple filename collision avoidance (prepend timestamp)
	fileName := fmt.Sprintf("%d_%s", os.Getpid(), file.Filename)
	filePath := filepath.Join(uploadDir, fileName)

	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		_ = ctx.Error(apperrors.Internal("failed to save file"))
		ctx.Abort()
		return
	}

	attachment, serviceErr := c.noteService.AddAttachment(ctx.Request.Context(), noteID, userID, file.Filename, filePath)
	if serviceErr != nil {
		// Clean up file if DB fails
		os.Remove(filePath)
		_ = ctx.Error(serviceErr)
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": attachment})
}
