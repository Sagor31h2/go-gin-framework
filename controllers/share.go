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
// ShareNote godoc
// @Summary Share a note with another user
// @Description Grant access to a note for another user
// @Tags shares
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Note ID"
// @Param body body object true "Target User ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /notes/{id}/share [post]
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
// RevokeShare godoc
// @Summary Revoke a note share
// @Description Remove access to a note for another user
// @Tags shares
// @Security BearerAuth
// @Param id path string true "Note ID"
// @Param uid path string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /notes/{id}/share/{uid} [delete]
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
// ListShares godoc
// @Summary List all shares for a note
// @Description Get a list of users who have access to the note
// @Tags shares
// @Security BearerAuth
// @Produce json
// @Param id path string true "Note ID"
// @Success 200 {object} map[string][]models.NoteShare
// @Failure 404 {object} map[string]string
// @Router /notes/{id}/shares [get]
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
