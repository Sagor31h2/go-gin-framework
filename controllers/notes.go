package controllers

import (
	"net/http"

	models "gin-test/internal/models"
	"gin-test/services"

	"github.com/gin-gonic/gin"
)

type NoteController struct {
	noteService services.NoteService
}

func NewNoteController(router *gin.Engine, noteService services.NoteService) *NoteController {
	c := &NoteController{
		noteService: noteService,
	}

	notes := router.Group("/notes")
	{
		notes.GET("/", c.GetNotes)
		notes.POST("/", c.CreateNote)
	}

	return c
}

func (c *NoteController) GetNotes(ctx *gin.Context) {
	notes, err := c.noteService.GetNotes(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch notes"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": notes,
	})
}

func (c *NoteController) CreateNote(ctx *gin.Context) {
	var note models.Note
	if err := ctx.ShouldBindJSON(&note); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.noteService.CreateNote(ctx.Request.Context(), &note); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create note"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data": note,
	})
}
