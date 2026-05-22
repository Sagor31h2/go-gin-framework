package controllers

import (
	"net/http"
	"strconv"

	apperrors "gin-test/internal/errors"
	"gin-test/services"

	"github.com/gin-gonic/gin"
)

type ShareController struct {
	noteService services.NoteService
}

func NewShareController(noteService services.NoteService) *ShareController {
	return &ShareController{noteService: noteService}
}

// POST /api/v1/notes/:id/share
func (c *ShareController) ShareNote(ctx *gin.Context) {
	ownerID := ctx.MustGet("user_id").(uint)
	noteID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		_ = ctx.Error(apperrors.BadRequest("invalid note id"))
		ctx.Abort()
		return
	}
	var body struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		_ = ctx.Error(apperrors.BadRequest(err.Error()))
		ctx.Abort()
		return
	}
	if err := c.noteService.ShareNote(ctx.Request.Context(), uint(noteID), ownerID, body.UserID); err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "note shared"})
}

// DELETE /api/v1/notes/:id/share/:uid
func (c *ShareController) RevokeShare(ctx *gin.Context) {
	ownerID := ctx.MustGet("user_id").(uint)
	noteID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		_ = ctx.Error(apperrors.BadRequest("invalid note id"))
		ctx.Abort()
		return
	}
	targetUID, err := strconv.ParseUint(ctx.Param("uid"), 10, 64)
	if err != nil {
		_ = ctx.Error(apperrors.BadRequest("invalid user id"))
		ctx.Abort()
		return
	}
	if err := c.noteService.RevokeShare(ctx.Request.Context(), uint(noteID), ownerID, uint(targetUID)); err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "share revoked"})
}

// GET /api/v1/notes/:id/shares
func (c *ShareController) ListShares(ctx *gin.Context) {
	ownerID := ctx.MustGet("user_id").(uint)
	noteID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		_ = ctx.Error(apperrors.BadRequest("invalid note id"))
		ctx.Abort()
		return
	}
	shares, err := c.noteService.ListShares(ctx.Request.Context(), uint(noteID), ownerID)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": shares})
}
