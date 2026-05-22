package controllers

import (
	"fmt"
	"net/http"

	models "gin-test/internal/models"
	"gin-test/services"

	"github.com/gin-gonic/gin"
)

type NoteController struct {
	noteService services.NoteService
}

func NewNoteController(noteService services.NoteService) *NoteController {
	return &NoteController{
		noteService: noteService,
	}
}

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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch notes"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  notes,
		"limit": limit,
		"page":  page,
	})
}

func (c *NoteController) GetNoteByID(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(uint)
	id := ctx.Param("id")
	note, err := c.noteService.GetNoteByID(ctx.Request.Context(), id, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": note,
	})
}

func (c *NoteController) CreateNote(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(uint)
	var note models.Note
	if err := ctx.ShouldBindJSON(&note); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	note.UserID = userID

	if err := c.noteService.CreateNote(ctx.Request.Context(), &note); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create note"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data": note,
	})
}

func (c *NoteController) UpdateNote(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(uint)
	id := ctx.Param("id")
	var note models.Note
	if err := ctx.ShouldBindJSON(&note); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.noteService.UpdateNote(ctx.Request.Context(), id, userID, &note); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update note"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "note updated successfully",
	})
}

func (c *NoteController) DeleteNote(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(uint)
	id := ctx.Param("id")
	if err := c.noteService.DeleteNote(ctx.Request.Context(), id, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete note"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "note deleted successfully",
	})
}
