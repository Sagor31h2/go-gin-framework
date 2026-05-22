package controllers

import (
	"fmt"
	"net/http"

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

func (c *NoteController) DeleteNote(ctx *gin.Context) {
	userID := ctx.MustGet("user_id").(uint)
	if err := c.noteService.DeleteNote(ctx.Request.Context(), ctx.Param("id"), userID); err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "note deleted successfully"})
}
