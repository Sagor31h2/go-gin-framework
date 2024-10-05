package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type NoteControllers struct{}

func (n *NoteControllers) InitNotesControllersRoutes(router *gin.Engine) {
	notes := router.Group("/notes")

	notes.GET("/", n.GetNotes())
	notes.POST("/", n.CreateNote())
}

func (n *NoteControllers) GetNotes() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "message:notes")
	}
}

func (n *NoteControllers) CreateNote() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "message:note created")
	}
}
