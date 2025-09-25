package controllers

import (
	"net/http"

	services "gin-test/Services"

	"github.com/gin-gonic/gin"
)

type NoteController struct {
	noteService services.NotesService
}

func (n *NoteController) InitNotesControllerRoutes(router *gin.Engine, noteService services.NotesService) {
	n.noteService = noteService
	notes := router.Group("/notes")

	notes.GET("/", n.GetNotes)
	notes.POST("/", n.CreateNote)
}

func (n *NoteController) GetNotes(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": n.noteService.GetNotes(),
	})
}

func (n *NoteController) CreateNote(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": n.noteService.CreateNotes(),
	})
}
